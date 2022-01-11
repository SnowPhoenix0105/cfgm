package tree2obj

import (
	"github.com/SnowPhoenix0105/cfgm/internal/tree"
	"reflect"
)

func Refill(
	root *tree.Node,
	obj interface{},
	buildTime tree.ModifyTime,
	currentTime tree.ModifyTime) {
	env := refillEnv{
		walker:    tree.WriteFrom(root, currentTime),
		buildTime: buildTime,
	}
	env.refill(reflect.Indirect(reflect.ValueOf(obj)))
}
