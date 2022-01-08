package tree

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFullNodeFlag(t *testing.T) {
	assert.Equal(t, fullNodeFlag(0), flagEmpty)
	assert.Equal(t, fullNodeFlag(1), flagInt)
	assert.Equal(t, fullNodeFlag(2), flagFloat)
	assert.Equal(t, fullNodeFlag(4), flagBool)
}

func TestFullNodeFlag_Add(t *testing.T) {
	flags := flagEmpty
	flags.add(flagBool)
	assert.Equal(t, flagBool, flags)
	flags.add(flagString)
	assert.Equal(t, flagBool|flagString, flags)
}

func TestFullNodeFlag_Delete(t *testing.T) {
	flags := flagString | flagFloat | flagObj
	flags.delete(flagString)
	assert.Equal(t, flagFloat|flagObj, flags)
	flags.delete(flagFloat)
	assert.Equal(t, flagObj, flags)
	flags.delete(flagObj)
	assert.Equal(t, flagEmpty, flags)
}

func TestFullNode_ImplementInnerNode(t *testing.T) {
	var _ InnerNode = newFullNode()
}

func TestFullNode_SetInt(t *testing.T) {
	var node NodeReadWriter = &Node{newFullNode()}
	assert.False(t, node.Has(NodeKeyInt))
	node.SetInt(123)
	assert.True(t, node.Has(NodeKeyInt))
	assert.Equal(t, int64(123), node.Int())
}

func TestFullNode_SetFloat(t *testing.T) {
	var node NodeReadWriter = &Node{newFullNode()}
	assert.False(t, node.Has(NodeKeyFloat))
	node.SetFloat(1.23)
	assert.True(t, node.Has(NodeKeyFloat))
	assert.Equal(t, 1.23, node.Float())
}

func TestFullNode_SetBool(t *testing.T) {
	var node NodeReadWriter = &Node{newFullNode()}
	assert.False(t, node.Has(NodeKeyBool))
	node.SetBool(true)
	assert.True(t, node.Has(NodeKeyBool))
	assert.True(t, node.Bool())
}

func TestFullNode_SetString(t *testing.T) {
	var node NodeReadWriter = &Node{newFullNode()}
	assert.False(t, node.Has(NodeKeyString))
	node.SetString("abc")
	assert.True(t, node.Has(NodeKeyString))
	assert.Equal(t, "abc", node.String())
}

func TestFullNode_SetObj(t *testing.T) {
	var node NodeReadWriter = &Node{newFullNode()}
	assert.False(t, node.Has(NodeKeyObj))
	obj := NodeObj{
		"a": NewNode(),
		"b": NewNode(),
	}
	node.SetObj(obj)
	assert.True(t, node.Has(NodeKeyObj))
	assert.Equal(t, obj, node.Obj())
}

func TestFullNode_SetList(t *testing.T) {
	var node NodeReadWriter = &Node{newFullNode()}
	assert.False(t, node.Has(NodeKeyList))
	list := NodeList{
		NewNode(),
		NewNode(),
	}
	node.SetList(list)
	assert.True(t, node.Has(NodeKeyList))
	assert.Equal(t, list, node.List())
}

func TestFullNode_SetObjPrototype(t *testing.T) {
	var node NodeReadWriter = &Node{newFullNode()}
	prototype := NewNode()
	assert.False(t, node.Has(NodeKeyObjPrototype))
	node.SetObjPrototype(prototype)
	assert.True(t, node.Has(NodeKeyObjPrototype))
	assert.Equal(t, prototype, node.ObjPrototype())
}

func TestFullNode_SetListPrototype(t *testing.T) {
	var node NodeReadWriter = &Node{newFullNode()}
	prototype := NewNode()
	assert.False(t, node.Has(NodeKeyListPrototype))
	node.SetListPrototype(prototype)
	assert.True(t, node.Has(NodeKeyListPrototype))
	assert.Equal(t, prototype, node.ListPrototype())
}

func TestFullNode_SetModifyTimeFor(t *testing.T) {
	var node NodeReadWriter = &Node{newFullNode()}

	node.SetModifyTimeFor(NodeKeyObj, 1)
	assert.Equal(t, ModifyTime(1), node.ModifyTimeFor(NodeKeyObj))
}

func TestFullNode_Delete(t *testing.T) {
	var node NodeReadWriter = &Node{newFullNode()}
	node.SetBool(true)
	node.SetInt(123)
	node.SetObjPrototype(NewNode())
	assert.True(t, node.Has(NodeKeyBool))
	assert.True(t, node.Has(NodeKeyInt))
	assert.True(t, node.Has(NodeKeyObjPrototype))
	node.Delete(NodeKeyBool)
	assert.False(t, node.Has(NodeKeyBool))
	assert.True(t, node.Has(NodeKeyInt))
	assert.True(t, node.Has(NodeKeyObjPrototype))
}
