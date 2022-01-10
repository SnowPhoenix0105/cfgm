package tree2json

import (
	"github.com/SnowPhoenix0105/cfgm/internal/tree"
	"io"
	"strings"
)

func DumpToString(root *tree.Node) string {
	writer := stringBuilderWriter{
		builder:        strings.Builder{},
		level:          0,
		commentLevel:   make(map[int]emptyType),
		endLineNeedNew: false,
	}
	env := dumpEnv{
		json:   &writer,
		walker: tree.ReadFrom(root),
	}
	env.dump()
	return writer.builder.String()
}

func DumpToWriter(root *tree.Node, writer io.Writer) (int, error) {
	str := DumpToString(root)
	return writer.Write([]byte(str))
}
