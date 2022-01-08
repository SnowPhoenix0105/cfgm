package tree

import (
	"errors"
	"fmt"
)

type fullNodeFlag int32

const (
	flagEmpty fullNodeFlag = 0
	flagInt   fullNodeFlag = 1 << (iota - 1)
	flagFloat
	flagBool
	flagString
	flagObj
	flagList
	flagObjPrototype
	flagListPrototype
)

func (flag fullNodeFlag) has(target fullNodeFlag) bool {
	if flag == flagEmpty {
		return false
	}
	return (flag & target) == target
}

func (flag *fullNodeFlag) add(target fullNodeFlag) {
	*flag |= target
}

func (flag *fullNodeFlag) delete(target fullNodeFlag) {
	*flag &= ^target
}

func keyToFlag(key NodeKey) fullNodeFlag {
	switch key {
	case NodeKeyInt:
		return flagInt
	case NodeKeyFloat:
		return flagFloat
	case NodeKeyBool:
		return flagBool
	case NodeKeyString:
		return flagString
	case NodeKeyObj:
		return flagObj
	case NodeKeyList:
		return flagList
	case NodeKeyObjPrototype:
		return flagObjPrototype
	case NodeKeyListPrototype:
		return flagListPrototype
	default:
		return flagEmpty
	}
}

type fullNode struct {
	flags      fullNodeFlag
	modifyTime map[NodeKey]ModifyTime

	intValue    int64
	floatValue  float64
	boolValue   bool
	stringValue string
	objValue    NodeObj
	listValue   NodeList

	objPrototype  *Node
	listPrototype *Node
}

func newFullNode() *fullNode {
	ret := new(fullNode)
	ret.init()
	return ret
}

func makeFullNode() fullNode {
	ret := fullNode{}
	ret.init()
	return ret
}

func (node *fullNode) init() {
	node.flags = flagEmpty
	node.modifyTime = make(map[NodeKey]ModifyTime)
}

func (node *fullNode) Has(key NodeKey) bool {
	flag := keyToFlag(key)
	if flag == flagEmpty {
		return false
	}
	return node.flags.has(flag)
}

func (node *fullNode) ModifyTimeFor(key NodeKey) ModifyTime {
	time, ok := node.modifyTime[key]
	if !ok {
		panic(errors.New(fmt.Sprintf("modify time for %s not exist", key.String())))
	}
	return time
}

func (node *fullNode) SetModifyTimeFor(key NodeKey, time ModifyTime) {
	node.modifyTime[key] = time
}

func (node *fullNode) Delete(key NodeKey) {
	flag := keyToFlag(key)
	if flag == flagEmpty {
		return
	}
	node.flags.delete(flag)
}

func (node *fullNode) Int() int64 {
	return node.intValue
}

func (node *fullNode) Float() float64 {
	return node.floatValue
}

func (node *fullNode) Bool() bool {
	return node.boolValue
}

func (node *fullNode) String() string {
	return node.stringValue
}

func (node *fullNode) Obj() NodeObj {
	return node.objValue
}

func (node *fullNode) List() NodeList {
	return node.listValue
}

func (node *fullNode) ObjPrototype() *Node {
	return node.objPrototype
}

func (node *fullNode) ListPrototype() *Node {
	return node.listPrototype
}

func (node *fullNode) SetInt(value int64) InnerNode {
	node.intValue = value
	node.flags.add(flagInt)
	return node
}

func (node *fullNode) SetFloat(value float64) InnerNode {
	node.floatValue = value
	node.flags.add(flagFloat)
	return node
}

func (node *fullNode) SetBool(value bool) InnerNode {
	node.boolValue = value
	node.flags.add(flagBool)
	return node
}

func (node *fullNode) SetString(value string) InnerNode {
	node.stringValue = value
	node.flags.add(flagString)
	return node
}

func (node *fullNode) SetObj(value NodeObj) InnerNode {
	node.objValue = value
	node.flags.add(flagObj)
	return node
}

func (node *fullNode) SetList(value NodeList) InnerNode {
	node.listValue = value
	node.flags.add(flagList)
	return node
}

func (node *fullNode) SetObjPrototype(value *Node) InnerNode {
	node.objPrototype = value
	node.flags.add(flagObjPrototype)
	return node
}

func (node *fullNode) SetListPrototype(value *Node) InnerNode {
	node.listPrototype = value
	node.flags.add(flagListPrototype)
	return node
}
