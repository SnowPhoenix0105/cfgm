package json2tree

import "github.com/SnowPhoenix0105/cfgm/internal/tree"

func setNull(walker tree.Walker) {
	handler := setNullDistributeHandler{walker: walker}
	tree.DistributeOnWalker(walker, &handler)
}

type setNullDistributeHandler struct {
	walker tree.Walker
}

func setNullIfNullable(walker tree.Walker, key tree.NodeKey) {
	if walker.NullableFor(key) {
		walker.SetNullFor(key, true)
	}
}

func (s *setNullDistributeHandler) HandleInt() {
	setNullIfNullable(s.walker, tree.NodeKeyInt)
}

func (s *setNullDistributeHandler) HandleFloat() {
	setNullIfNullable(s.walker, tree.NodeKeyFloat)
}

func (s *setNullDistributeHandler) HandleBool() {
	setNullIfNullable(s.walker, tree.NodeKeyBool)
}

func (s *setNullDistributeHandler) HandleString() {
	setNullIfNullable(s.walker, tree.NodeKeyString)
}

func (s *setNullDistributeHandler) HandleObj() {
	setNullIfNullable(s.walker, tree.NodeKeyObj)
}

func (s *setNullDistributeHandler) HandleList() {
	setNullIfNullable(s.walker, tree.NodeKeyList)
}
