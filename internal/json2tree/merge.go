package json2tree

import (
	"github.com/SnowPhoenix0105/cfgm/internal/tree"
	"io"
)

func Merge(root *tree.Node, reader io.RuneReader, time tree.ModifyTime) error {
	par := parser{
		walker: tree.WriteFrom(root, time),
	}
	err := par.Reset(reader)
	if err != nil {
		return err
	}
	return par.parseNode()
}
