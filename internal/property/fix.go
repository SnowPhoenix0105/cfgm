package property

import (
	"github.com/SnowPhoenix0105/cfgm/internal/tree"
	"strconv"
)

type fixEnv struct {
	walker tree.Walker
}

func (env *fixEnv) getListLength() int {
	if !env.walker.Has(tree.NodeKeyList) {
		return 0
	}
	return len(env.walker.List())
}

func (env *fixEnv) assign(value string) {
	length := len(value)
	if length == 0 {
		if env.walker.Has(tree.NodeKeyBool) {
			env.walker.SetBool(true)
		}
		return
	}
	prefix := value[0]
	if prefix == '0' && len(value) > 2 {
		base := 8
		beg := 1
		switch value[1] {
		case 'x':
			base = 16
			beg = 2
		case 'b':
			base = 2
			beg = 2
		}
		i, err := strconv.ParseInt(value[beg:length-1], base, 64)
		if err == nil {
			env.walker.SetInt(i)
			return
		}
	}
	if '0' <= prefix && prefix <= '9' {
		i, err := strconv.ParseInt(value, 10, 64)
		if err == nil {
			env.walker.SetInt(i)
			return
		}
		f, err := strconv.ParseFloat(value, 64)
		if err == nil {
			env.walker.SetFloat(f)
			return
		}
	}
	if prefix == '"' && value[length-1] == '"' {
		env.walker.SetString(value[1 : length-1])
		return
	}
	if value == "true" {
		env.walker.SetBool(true)
		return
	}
	if value == "false" {
		env.walker.SetBool(false)
		return
	}
	env.walker.SetString(value)
}

func listIndex(base int, name string) int {
	if len(name) > 2 && name[0] == '[' && name[len(name)-1] == ']' {
		i, err := strconv.ParseInt(name[1:len(name)-1], 10, 32)
		if err != nil {
			return -1
		}
		index := int(i)
		if index < 0 || name[1] == '+' {
			index += base
		}
		return index
	}
	return -1
}

func (env *fixEnv) fixNode(ptr *node) {
	env.assign(ptr.value)
	baseLength := env.getListLength()
	for k, v := range ptr.sub {
		index := listIndex(baseLength, k)
		if index >= 0 {
			env.walker.EnterList(index)
			env.fixNode(v)
			env.walker.Exit()
			continue
		}
		env.walker.EnterObj(k)
		env.fixNode(v)
		env.walker.Exit()
	}
}
