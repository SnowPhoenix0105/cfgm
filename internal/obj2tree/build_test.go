package obj2tree

import (
	"github.com/SnowPhoenix0105/cfgm/internal/tree"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildFrom_Simple(t *testing.T) {
	obj := map[string]int{
		"A": 1,
		"B": 2,
	}

	actual, err := BuildFrom(obj, 1)

	expect := tree.NewNode()
	handler := tree.WriteFrom(expect, 1)
	handler.SetClearWhenEnterFor(tree.NodeKeyObj, true)
	handler.EnterObjPrototype()
	handler.SetInt(0)
	handler.Exit()
	handler.EnterObj("A")
	handler.SetInt(1)
	handler.Exit()
	handler.EnterObj("B")
	handler.SetInt(2)
	handler.Exit()

	assert.Nil(t, err)
	assert.True(t, tree.Equals(expect, actual))
}

func TestBuildFrom_ObjPrototype(t *testing.T) {
	obj := map[string]map[string]int{
		"__prototype__": map[string]int{
			"__prototype__": 1,
		},
	}

	actual, err := BuildFrom(obj, 1)

	expect := tree.NewNode()
	handler := tree.WriteFrom(expect, 1)
	handler.SetClearWhenEnterFor(tree.NodeKeyObj, true)
	handler.EnterObjPrototype()
	handler.SetClearWhenEnterFor(tree.NodeKeyObj, true)
	handler.EnterObjPrototype()
	handler.SetInt(1)
	handler.Exit()
	handler.Exit()

	assert.Nil(t, err)
	assert.True(t, tree.Equals(expect, actual))
	assert.Empty(t, obj)
	assert.Equal(t, 0, len(obj))
}

func TestBuildFrom_ListPrototype(t *testing.T) {
	obj := map[string][]int{
		"A": []int{
			1,
			2,
		}[:1],
	}

	actual, err := BuildFrom(obj, 1)

	expect := tree.NewNode()
	handler := tree.WriteFrom(expect, 1)
	handler.SetClearWhenEnterFor(tree.NodeKeyObj, true)
	handler.EnterObjPrototype()
	handler.SetClearWhenEnterFor(tree.NodeKeyList, true)
	handler.EnterListPrototype()
	handler.SetInt(0)
	handler.Exit()
	handler.Exit()
	handler.EnterObj("A")
	handler.SetClearWhenEnterFor(tree.NodeKeyList, true)
	handler.EnterListPrototype()
	handler.SetInt(1)
	handler.Exit()
	handler.Exit()

	assert.Nil(t, err)
	assert.True(t, tree.Equals(expect, actual))
	assert.Equal(t, 1, len(obj))
	assert.Equal(t, 0, len(obj["A"]))
}
