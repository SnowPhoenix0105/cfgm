package tree

type DistributeHandler interface {
	HandleInt()
	HandleFloat()
	HandleBool()
	HandleString()
	HandleObj()
	HandleList()
}

func Distribute(node *Node, handler DistributeHandler) {
	if node.Has(NodeKeyObjPrototype) || node.Has(NodeKeyObj) {
		handler.HandleObj()
	} else if node.Has(NodeKeyString) {
		handler.HandleString()
	} else if node.Has(NodeKeyInt) {
		handler.HandleInt()
	} else if node.Has(NodeKeyFloat) {
		handler.HandleFloat()
	} else if node.Has(NodeKeyBool) {
		handler.HandleBool()
	} else if node.Has(NodeKeyListPrototype) || node.Has(NodeKeyList) {
		handler.HandleList()
	}
}

func DistributeOnWalker(walker ReadonlyWalker, handler DistributeHandler) {
	if walker.Has(NodeKeyObjPrototype) || walker.Has(NodeKeyObj) {
		handler.HandleObj()
	} else if walker.Has(NodeKeyString) {
		handler.HandleString()
	} else if walker.Has(NodeKeyInt) {
		handler.HandleInt()
	} else if walker.Has(NodeKeyFloat) {
		handler.HandleFloat()
	} else if walker.Has(NodeKeyBool) {
		handler.HandleBool()
	} else if walker.Has(NodeKeyListPrototype) || walker.Has(NodeKeyList) {
		handler.HandleList()
	}
}
