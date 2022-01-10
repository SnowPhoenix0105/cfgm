package json2tree

import (
	"github.com/SnowPhoenix0105/cfgm/internal/obj2tree"
	"github.com/SnowPhoenix0105/cfgm/internal/tree2json"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestMapPrototype(t *testing.T) {
	type Tuple struct {
		A int
		B int
	}
	type Class struct {
		CoverMap      map[string]Tuple
		MergeMap      map[string]*Tuple
		CoverProtoMap map[string]Tuple
		MergeProtoMap map[string]*Tuple
	}
	obj := Class{
		CoverMap: map[string]Tuple{
			"Key1": {
				A: 1,
				B: 2,
			},
		},
		MergeMap: map[string]*Tuple{
			"Key1": {
				A: 1,
				B: 2,
			},
		},
		CoverProtoMap: map[string]Tuple{
			"__prototype__": {
				A: 3,
				B: 4,
			},
			"Key1": {
				A: 1,
				B: 2,
			},
		},
		MergeProtoMap: map[string]*Tuple{
			"__prototype__": {
				A: 3,
				B: 4,
			},
			"Key1": {
				A: 1,
				B: 2,
			},
		},
	}
	json := `
{
	"CoverMap": {
		"Key1": {
			"A": 5
		},
		"Key2": {
			"A": 6
		}
	},
	"MergeMap": {
		"Key1": {
			"A": 5
		},
		"Key2": {
			"A": 6
		}
	},
	"CoverProtoMap": {
		"Key1": {
			"A": 5
		},
		"Key2": {
			"A": 6
		}
	},
	"MergeProtoMap": {
		"Key1": {
			"A": 5
		},
		"Key2": {
			"A": 6
		}
	},
}`
	root, err := obj2tree.BuildFrom(&obj, 1)
	assert.Nil(t, err)
	reader := strings.NewReader(json)
	err = Merge(root, reader, 2)
	assert.Nil(t, err)
	str := tree2json.DumpToString(root)
	t.Log(str)
}
