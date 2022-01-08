package obj2tree

import (
	"github.com/SnowPhoenix0105/cfgm/internal/tree"
	"reflect"
)

func BuildFrom(obj interface{}, time tree.ModifyTime) (*tree.Node, error) {
	root := tree.NewNode()
	env := buildEnv{
		Walker:       tree.WriteFrom(root, time),
		DescTag:      "desc",
		PrototypeKey: "__prototype__",
	}
	err := env.buildFrom(reflect.ValueOf(obj), kvProperty{false, false})
	return root, err
}
