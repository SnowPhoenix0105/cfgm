package tree

func Equals(left, right *Node) bool {
	lRaw := left.Raw
	rRaw := right.Raw
	switch lNode := lRaw.(type) {
	case *fullNode:
		rNode, ok := rRaw.(*fullNode)
		if !ok {
			return false
		}
		return fullNodeEquals(lNode, rNode)
	}
	panic("not implement")
}

func fullNodeEquals(left, right *fullNode) bool {
	if left.flags != right.flags {
		return false
	}
	if left.modifyTime != right.modifyTime {
		return false
	}
	if left.Has(NodeKeyDesc) && left.descValue != right.descValue {
		return false
	}
	if left.Has(NodeKeyInt) && left.intValue != right.intValue {
		return false
	}
	if left.Has(NodeKeyFloat) && left.floatValue != right.floatValue {
		return false
	}
	if left.Has(NodeKeyBool) && left.boolValue != right.boolValue {
		return false
	}
	if left.Has(NodeKeyString) && left.stringValue != right.stringValue {
		return false
	}
	if left.Has(NodeKeyObjPrototype) && !Equals(left.objPrototype, right.objPrototype) {
		return false
	}
	if left.Has(NodeKeyListPrototype) && !Equals(left.listPrototype, right.listPrototype) {
		return false
	}
	if left.Has(NodeKeyList) {
		if len(left.listValue) != len(right.listValue) {
			return false
		}
		for i, e := range left.listValue {
			if !Equals(e, right.listValue[i]) {
				return false
			}
		}
	}
	if left.Has(NodeKeyObj) {
		if len(left.objValue) != len(right.objValue) {
			return false
		}
		for k, v := range left.objValue {
			v2, ok := right.objValue[k]
			if !ok {
				return false
			}
			if !Equals(v, v2) {
				return false
			}
		}
	}
	return true
}
