package obj2tree

import (
	"github.com/SnowPhoenix0105/cfgm/internal/tree"
	"reflect"
)

type buildEnv struct {
	DescTag             string
	PrototypeKey        string
	Walker              tree.Walker
	prototypeKeyReflect reflect.Value
}

type kvProperty struct {
	isNull   bool
	nullable bool
}

func (env *buildEnv) buildFrom(obj reflect.Value, property kvProperty) error {
	return env.distribute(obj, property)
}

func (env *buildEnv) distribute(obj reflect.Value, property kvProperty) error {
	switch obj.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return env.buildFromInt(obj, property)
	case reflect.Float32, reflect.Float64:
		return env.buildFromFloat(obj, property)
	case reflect.Bool:
		return env.buildFromBool(obj, property)
	case reflect.String:
		return env.buildFromString(obj, property)
	case reflect.Struct:
		return env.buildFromStruct(obj, property)
	case reflect.Slice:
		return env.buildFromSlice(obj, property)
	case reflect.Map:
		return env.buildFromMap(obj, property)
	}
	// TODO error
	panic("not implement")
}

func (env *buildEnv) buildFromInt(obj reflect.Value, property kvProperty) error {
	env.Walker.SetInt(obj.Int())
	env.Walker.SetNullFor(tree.NodeKeyInt, property.isNull)
	env.Walker.SetNullableFor(tree.NodeKeyInt, property.nullable)
	return nil
}

func (env *buildEnv) buildFromFloat(obj reflect.Value, property kvProperty) error {
	env.Walker.SetFloat(obj.Float())
	env.Walker.SetNullFor(tree.NodeKeyFloat, property.isNull)
	env.Walker.SetNullableFor(tree.NodeKeyFloat, property.nullable)
	return nil
}

func (env *buildEnv) buildFromBool(obj reflect.Value, property kvProperty) error {
	env.Walker.SetBool(obj.Bool())
	env.Walker.SetNullFor(tree.NodeKeyBool, property.isNull)
	env.Walker.SetNullableFor(tree.NodeKeyBool, property.nullable)
	return nil
}

func (env *buildEnv) buildFromString(obj reflect.Value, property kvProperty) error {
	env.Walker.SetString(obj.String())
	env.Walker.SetNullFor(tree.NodeKeyString, property.isNull)
	env.Walker.SetNullableFor(tree.NodeKeyString, property.nullable)
	return nil
}

func (env *buildEnv) buildFromStruct(obj reflect.Value, property kvProperty) error {
	typ := obj.Type()
	numField := obj.NumField()
	for i := 0; i < numField; i++ {
		fieldType := typ.Field(i)
		fieldValue := obj.Field(i)
		err := env.buildFromField(fieldValue, fieldType)
		if err != nil {
			return err
		}
	}

	env.Walker.SetNullFor(tree.NodeKeyObj, property.isNull)
	env.Walker.SetNullableFor(tree.NodeKeyObj, property.nullable)
	return nil
}

func (env *buildEnv) unwrapPointerForField(ptr reflect.Value) (elem reflect.Value, property kvProperty) {
	if ptr.Kind() != reflect.Ptr {
		return ptr, kvProperty{false, false}
	}
	elem1 := ptr.Elem()
	if !elem1.IsValid() {
		// nil ptr
		return reflect.Zero(ptr.Type().Elem()), kvProperty{true, true}
	}
	if elem1.Kind() == reflect.Ptr {
		// pointer to pointer
		elem2 := elem1.Elem()
		if !elem2.IsValid() {
			return reflect.Zero(elem1.Type().Elem()), kvProperty{true, true}
		}
		return elem2, kvProperty{true, true}
	}
	return elem1, kvProperty{false, true}
}

func (env *buildEnv) buildFromField(obj reflect.Value, field reflect.StructField) error {
	env.Walker.EnterObj(field.Name)
	defer env.Walker.Exit()

	elem, property := env.unwrapPointerForField(obj)
	// TODO check
	err := env.buildFrom(elem, property)
	if err != nil {
		return err
	}

	desc := field.Tag.Get(env.DescTag)
	if len(desc) != 0 {
		env.Walker.SetDesc(desc)
	}
	return nil
}

func (env *buildEnv) unwrapPointerForElem(ptr reflect.Value) (elem reflect.Value, property kvProperty) {
	if ptr.Kind() != reflect.Ptr {
		return ptr, kvProperty{false, false}
	}
	elem1 := ptr.Elem()
	if !elem1.IsValid() {
		// nil ptr
		return reflect.Zero(ptr.Type().Elem()), kvProperty{true, true}
	}
	return elem1, kvProperty{false, true}
}

func (env *buildEnv) buildFromSlice(obj reflect.Value, property kvProperty) error {
	isPtr := obj.Type().Elem().Kind() == reflect.Ptr
	// TODO check
	env.Walker.SetClearWhenEnterFor(tree.NodeKeyList, !isPtr)

	length := obj.Len()
	capacity := obj.Cap()
	if capacity > length {
		length -= 1
		env.Walker.EnterListPrototype()
		elem := obj.Index(length)
		err := env.buildFrom(env.unwrapPointerForElem(elem))
		env.Walker.Exit()
		if err != nil {
			return err
		}
		obj.SetLen(length)
	} else {
		env.Walker.EnterListPrototype()
		typ := obj.Type().Elem()
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}
		err := env.buildFrom(reflect.Zero(typ), kvProperty{true, true})
		env.Walker.Exit()
		if err != nil {
			return err
		}
	}

	for i := 0; i < length; i++ {
		env.Walker.EnterList(i)
		elem := obj.Index(i)
		err := env.buildFrom(env.unwrapPointerForElem(elem))
		env.Walker.Exit()
		if err != nil {
			return err
		}
	}

	env.Walker.SetNullFor(tree.NodeKeyList, property.isNull)
	env.Walker.SetNullableFor(tree.NodeKeyList, property.nullable)
	return nil
}

func (env *buildEnv) buildFromMap(obj reflect.Value, property kvProperty) error {
	isPtr := obj.Type().Elem().Kind() == reflect.Ptr
	// TODO check
	env.Walker.SetClearWhenEnterFor(tree.NodeKeyList, !isPtr)

	if !env.prototypeKeyReflect.IsValid() {
		env.prototypeKeyReflect = reflect.ValueOf(env.PrototypeKey)
	}
	prototype := obj.MapIndex(env.prototypeKeyReflect)
	if prototype.IsValid() {
		env.Walker.EnterListPrototype()
		err := env.buildFrom(env.unwrapPointerForElem(prototype))
		env.Walker.Exit()
		if err != nil {
			return err
		}
		obj.SetMapIndex(env.prototypeKeyReflect, reflect.Value{})
	} else {
		env.Walker.EnterListPrototype()
		typ := obj.Type().Elem()
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}
		err := env.buildFrom(reflect.Zero(typ), kvProperty{true, true})
		env.Walker.Exit()
		if err != nil {
			return err
		}
	}

	iter := obj.MapRange()
	for iter.Next() {
		env.Walker.EnterObj(iter.Key().String())
		err := env.buildFrom(env.unwrapPointerForElem(iter.Value()))
		env.Walker.Exit()
		if err != nil {
			return err
		}
	}

	env.Walker.SetNullFor(tree.NodeKeyObj, property.isNull)
	env.Walker.SetNullableFor(tree.NodeKeyObj, property.nullable)
	return nil
}
