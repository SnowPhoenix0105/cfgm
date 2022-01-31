package tree2obj

import (
	"fmt"
	"github.com/SnowPhoenix0105/cfgm/internal/tree"
	"reflect"
)

type refillEnv struct {
	walker    tree.ReadonlyWalker
	buildTime tree.ModifyTime
}

func (env *refillEnv) refill(obj reflect.Value) {
	if env.walker.ModifyTime() == env.buildTime {
		return
	}
	switch obj.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if env.walker.Has(tree.NodeKeyInt) {
			obj.SetInt(env.walker.Int())
		}
		return
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if env.walker.Has(tree.NodeKeyInt) {
			obj.SetUint(uint64(env.walker.Int()))
		}
		return
	case reflect.Float32, reflect.Float64:
		if env.walker.Has(tree.NodeKeyFloat) {
			obj.SetFloat(env.walker.Float())
		}
		return
	case reflect.Bool:
		if env.walker.Has(tree.NodeKeyBool) {
			obj.SetBool(env.walker.Bool())
		}
		return
	case reflect.String:
		if env.walker.Has(tree.NodeKeyString) {
			obj.SetString(env.walker.String())
		}
		return
	case reflect.Map:
		env.refillMap(obj)
		return
	case reflect.Struct:
		env.refillStruct(obj)
		return
	case reflect.Slice:
		env.refillSlice(obj)
		return
	}
	panic("not implement")
}

func (env *refillEnv) refillMap(obj reflect.Value) {
	if !env.walker.Has(tree.NodeKeyObj) {
		return
	}
	mapType := obj.Type()
	valueType := mapType.Elem()
	valueTypeIsPtr := valueType.Kind() == reflect.Ptr
	var elemType reflect.Type
	if valueTypeIsPtr {
		elemType = valueType.Elem()
	} else {
		elemType = valueType
		obj.Set(reflect.MakeMap(obj.Type()))
	}

	keys := env.walker.ObjKeys()
	for _, key := range keys {
		if !env.walker.TryEnterObj(key) {
			if DEBUG {
				panic("TryEnterObj() fail with key from ObjKeys()")
			}
			continue
		}
		keyReflect := reflect.ValueOf(key)
		if valueTypeIsPtr {
			ptr := obj.MapIndex(keyReflect)
			if ptr.IsValid() {
				elem := ptr.Elem()
				if elem.IsValid() {
					// merge by modify, and default value is provided
					env.refill(elem)
					env.walker.Exit()
					continue
				}
			}
		}
		ptr := reflect.New(elemType)
		elem := ptr.Elem()
		env.refill(elem)
		if valueTypeIsPtr {
			obj.SetMapIndex(keyReflect, ptr)
		} else {
			obj.SetMapIndex(keyReflect, elem)
		}
		env.walker.Exit()
	}
}

func (env *refillEnv) refillSlice(obj reflect.Value) {
	if !env.walker.Has(tree.NodeKeyList) {
		return
	}
	sliceType := obj.Type()
	valueType := sliceType.Elem()
	valueTypeIsPtr := valueType.Kind() == reflect.Ptr
	var elemType reflect.Type
	if valueTypeIsPtr {
		elemType = valueType.Elem()
	} else {
		elemType = valueType
	}

	length := env.walker.ListLen()
	var objLength int
	if !valueTypeIsPtr {
		obj.SetLen(0)
		objLength = 0
	} else {
		objLength = obj.Len()
	}
	for i := 0; i < length; i++ {
		if !env.walker.TryEnterList(i) {
			if DEBUG {
				panic("TryEnterList() fail with index less than ListLen()")
			}
			continue
		}
		if i < objLength {
			ptr := obj.Index(i)
			if ptr.IsValid() {
				elem := ptr.Elem()
				if elem.IsValid() {
					// merge by modify, and default value is provided
					env.refill(elem)
					env.walker.Exit()
					continue
				}
			}
		}
		ptr := reflect.New(elemType)
		elem := ptr.Elem()
		env.refill(elem)
		if valueTypeIsPtr {
			if i < objLength {
				obj.Index(i).Set(ptr)
			} else {
				obj.Set(reflect.Append(obj, ptr))
			}
		} else {
			obj.Set(reflect.Append(obj, elem))
		}
		env.walker.Exit()
	}
}

func (env *refillEnv) isNullFor(kind reflect.Kind) bool {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return env.walker.IsNullFor(tree.NodeKeyInt)
	case reflect.Float32, reflect.Float64:
		return env.walker.IsNullFor(tree.NodeKeyFloat)
	case reflect.Bool:
		return env.walker.IsNullFor(tree.NodeKeyBool)
	case reflect.String:
		return env.walker.IsNullFor(tree.NodeKeyString)
	case reflect.Map, reflect.Struct:
		return env.walker.IsNullFor(tree.NodeKeyObj)
	case reflect.Slice:
		return env.walker.IsNullFor(tree.NodeKeyList)
	}
	panic(fmt.Sprintf("Invalid Kind: %s", kind.String()))
}

func (env *refillEnv) refillPtrField(ptr reflect.Value) {
	if env.isNullFor(ptr.Type().Elem().Kind()) {
		ptr.Set(reflect.Zero(ptr.Type()))
		return
	}
	elem := ptr.Elem()
	if !elem.IsValid() {
		ptr.Set(reflect.New(ptr.Type().Elem()))
		elem = ptr.Elem()
	}
	env.refill(elem)
}

func (env *refillEnv) refillPtrPtrField(ptrptr reflect.Value) {
	if env.isNullFor(ptrptr.Type().Elem().Elem().Kind()) {
		ptrptr.Set(reflect.Zero(ptrptr.Type()))
		return
	}
	ptr := ptrptr.Elem()
	if !ptr.IsValid() {
		ptrptr.Set(reflect.New(ptrptr.Type().Elem()))
		ptr = ptrptr.Elem()
	}
	elem := ptr.Elem()
	if !elem.IsValid() {
		ptr.Set(reflect.New(ptr.Type().Elem()))
		elem = ptr.Elem()
	}
	env.refill(elem)
}

func (env *refillEnv) refillStruct(obj reflect.Value) {
	if !env.walker.Has(tree.NodeKeyObj) {
		return
	}

	structType := obj.Type()
	numField := obj.NumField()
	for i := 0; i < numField; i++ {
		field := structType.Field(i)
		elem := obj.Field(i)
		if !elem.CanSet() {
			continue
		}
		if !env.walker.TryEnterObj(field.Name) {
			continue
		}
		if elem.Kind() != reflect.Ptr {
			// simple
			env.refill(elem)
		} else if elem.Type().Elem().Kind() != reflect.Ptr {
			// ptr
			env.refillPtrField(elem)
		} else {
			// ptr to ptr
			env.refillPtrPtrField(elem)
		}
		env.walker.Exit()
	}
}
