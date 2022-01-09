package obj2tree

import (
	"github.com/SnowPhoenix0105/cfgm/internal/tree"
	"github.com/SnowPhoenix0105/deepcopy"
	"reflect"
)

func BuildFrom(obj interface{}, time tree.ModifyTime) (*tree.Node, error) {
	root := tree.NewNode()
	env := buildEnv{
		Walker:       tree.WriteFrom(root, time),
		DescTag:      "desc",
		PrototypeKey: "__prototype__",
		DeepCopy: deepcopy.WithOptions(&deepcopy.Options{
			IgnoreUnexploredFields: false,
		}),
	}
	err := env.buildFrom(reflect.ValueOf(obj), kvProperty{false, false})
	return root, err
}
