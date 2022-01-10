package tree2json

import (
	"github.com/SnowPhoenix0105/cfgm/internal/tree"
	"strconv"
)

type dumpEnv struct {
	json   jsonWriter
	walker tree.ReadonlyWalker
}

// <<<==== distribute begin ====>>>

func (env *dumpEnv) dump() {
	tree.DistributeOnWalker(env.walker, env)
}

func (env *dumpEnv) HandleInt() {
	env.json.WriteString(strconv.FormatInt(env.walker.Int(), 10))
}

func (env *dumpEnv) HandleFloat() {
	env.json.WriteString(strconv.FormatFloat(env.walker.Float(), 'f', 3, 64))
}

func (env *dumpEnv) HandleBool() {
	if env.walker.Bool() {
		env.json.WriteString("true")
	} else {
		env.json.WriteString("false")
	}
}

func (env *dumpEnv) HandleString() {
	env.dumpString(env.walker.String())
}

func (env *dumpEnv) HandleObj() {
	env.dumpObj()
}

func (env *dumpEnv) HandleList() {
	env.dumpList()
}

// <<----- distribute end ----->>

//func (env *dumpEnv) dump() {
//	if env.walker.Has(tree.NodeKeyObjPrototype) || env.walker.Has(tree.NodeKeyObj) {
//		env.dumpObj()
//	} else if env.walker.Has(tree.NodeKeyString) {
//		env.dumpString(env.walker.String())
//	} else if env.walker.Has(tree.NodeKeyInt) {
//		env.json.WriteString(strconv.FormatInt(env.walker.Int(), 10))
//	} else if env.walker.Has(tree.NodeKeyFloat) {
//		env.json.WriteString(strconv.FormatFloat(env.walker.Float(), 'f', 3, 64))
//	} else if env.walker.Has(tree.NodeKeyBool) {
//		if env.walker.Bool() {
//			env.json.WriteString("true")
//		} else {
//			env.json.WriteString("false")
//		}
//	} else if env.walker.Has(tree.NodeKeyListPrototype) || env.walker.Has(tree.NodeKeyList) {
//		env.dumpList()
//	}
//}

func (env *dumpEnv) dumpString(str string) {
	env.json.WriteRune('"')
	for _, char := range str {
		switch char {
		case '\t':
			env.json.WriteString("\\t")
		case '\n':
			env.json.WriteString("\\n")
		case '\b':
			env.json.WriteString("\\b")
		case '\v':
			env.json.WriteString("\\v")
		case '\r':
			env.json.WriteString("\\r")
		default:
			env.json.WriteRune(char)
		}
	}
	env.json.WriteRune('"')
}

func (env *dumpEnv) dumpObj() {
	env.json.WriteRune('{')
	env.json.Enter()
	env.json.EndLine()

	// prototype
	if env.walker.TryEnterObjPrototype() {
		env.json.StartComment()
		env.dumpString("Key")
		env.json.WriteRune(':')
		env.json.WriteSpace()
		env.dump()
		env.json.WriteRune(',')
		env.json.EndComment()
		env.walker.Exit()
	}

	// content
	keys := env.walker.ObjKeys()
	for i, key := range keys {
		if i != 0 {
			env.json.WriteRune(',')
		}
		env.json.EndLine()

		ok := env.walker.TryEnterObj(key)
		if DEBUG {
			if !ok {
				panic("TryEnterObj() cannot enter with key from walker.ObjKeys()")
			}
		}

		if env.walker.Has(tree.NodeKeyDesc) {
			env.json.CommentAndNewLine(env.walker.Desc())
		}

		env.dumpString(key)
		env.json.WriteRune(':')
		env.json.WriteSpace()

		env.dump()

		env.walker.Exit()
	}

	env.json.Exit()
	env.json.EndLine()
	env.json.WriteRune('}')
}

func (env *dumpEnv) dumpList() {
	env.json.WriteRune('[')
	env.json.Enter()
	env.json.EndLine()

	// prototype
	if env.walker.TryEnterListPrototype() {
		env.json.StartComment()
		env.dump()
		env.json.WriteRune(',')
		env.json.EndComment()
		env.walker.Exit()
	}

	// content
	length := env.walker.ListLen()
	for i := 0; i < length; i++ {
		if i != 0 {
			env.json.WriteRune(',')
		}
		env.json.EndLine()
		env.json.WriteSpace()

		ok := env.walker.TryEnterList(i)
		if DEBUG {
			if !ok {
				panic("TryEnterList() cannot enter with index less than walker.ListLen()")
			}
		}
		env.dump()
		env.walker.Exit()
	}

	env.json.Exit()
	env.json.EndLine()
	env.json.WriteRune(']')
}
