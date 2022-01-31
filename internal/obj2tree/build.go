package obj2tree

import (
	"github.com/SnowPhoenix0105/cfgm/internal/tree"
	"github.com/SnowPhoenix0105/deepcopy"
	"reflect"
)

func BuildFrom(obj interface{}, time tree.ModifyTime) (*tree.Node, error) {
	root := tree.NewNode()
	err := AppendTo(obj, tree.WriteFrom(root, time))
	return root, err
}

func AppendTo(obj interface{}, walker tree.Walker) error {
	env := buildEnv{
		Walker:       walker,
		DescTag:      "desc",
		PrototypeKey: "__prototype__",
		DeepCopy: deepcopy.WithOptions(&deepcopy.Options{
			IgnoreUnexploredFields: false,
		}),
	}

	return env.buildFrom(
		reflect.Indirect(reflect.ValueOf(obj)),
		kvProperty{false, false})
}
