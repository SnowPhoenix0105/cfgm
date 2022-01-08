package tree

type fullNodeFlag uint32

const (
	flagEmpty fullNodeFlag = 0

	flagHasDesc fullNodeFlag = 1 << (iota - 1)
	flagHasInt
	flagHasFloat
	flagHasBool
	flagHasString
	flagHasObj
	flagHasList

	flagIsNullDesc
	flagIsNullInt
	flagIsNullFloat
	flagIsNullBool
	flagIsNullString
	flagIsNullObj
	flagIsNullList

	flagNullableDesc
	flagNullableInt
	flagNullableFloat
	flagNullableBool
	flagNullableString
	flagNullableObj
	flagNullableList

	flagHasObjPrototype
	flagHasListPrototype

	flagClearObjWhenEnter
	flagClearListWhenEnter
)

func (flag fullNodeFlag) has(target fullNodeFlag) bool {
	if flag == flagEmpty {
		return false
	}
	return (flag & target) == target
}

func (flag *fullNodeFlag) add(target fullNodeFlag) {
	*flag |= target
}

func (flag *fullNodeFlag) delete(target fullNodeFlag) {
	*flag &= ^target
}

func keyToHasFlag(key NodeKey) fullNodeFlag {
	switch key {
	case NodeKeyDesc:
		return flagHasDesc
	case NodeKeyInt:
		return flagHasInt
	case NodeKeyFloat:
		return flagHasFloat
	case NodeKeyBool:
		return flagHasBool
	case NodeKeyString:
		return flagHasString
	case NodeKeyObj:
		return flagHasObj
	case NodeKeyList:
		return flagHasList
	case NodeKeyObjPrototype:
		return flagHasObjPrototype
	case NodeKeyListPrototype:
		return flagHasListPrototype
	default:
		return flagEmpty
	}
}

func keyToIsNullFlag(key NodeKey) fullNodeFlag {
	switch key {
	case NodeKeyDesc:
		return flagIsNullDesc
	case NodeKeyInt:
		return flagIsNullInt
	case NodeKeyFloat:
		return flagIsNullFloat
	case NodeKeyBool:
		return flagIsNullBool
	case NodeKeyString:
		return flagIsNullString
	case NodeKeyObj:
		return flagIsNullObj
	case NodeKeyList:
		return flagIsNullList
	default:
		return flagEmpty
	}
}

func keyToNullableFlag(key NodeKey) fullNodeFlag {
	switch key {
	case NodeKeyDesc:
		return flagNullableDesc
	case NodeKeyInt:
		return flagNullableInt
	case NodeKeyFloat:
		return flagNullableFloat
	case NodeKeyBool:
		return flagNullableBool
	case NodeKeyString:
		return flagNullableString
	case NodeKeyObj:
		return flagNullableObj
	case NodeKeyList:
		return flagNullableList
	default:
		return flagEmpty
	}
}

func keyToClearWhenEnterFlag(key NodeKey) fullNodeFlag {
	switch key {
	case NodeKeyObj:
		return flagClearObjWhenEnter
	case NodeKeyList:
		return flagClearListWhenEnter
	default:
		return flagEmpty
	}
}

type fullNode struct {
	flags      fullNodeFlag
	modifyTime ModifyTime

	descValue   string
	intValue    int64
	floatValue  float64
	boolValue   bool
	stringValue string
	objValue    NodeObj
	listValue   NodeList

	objPrototype  *Node
	listPrototype *Node
}

func newFullNode() *fullNode {
	ret := new(fullNode)
	return ret
}

func makeFullNode() fullNode {
	ret := fullNode{}
	return ret
}

// <<<==== readonly methods begin ====>>>

func (node *fullNode) Has(key NodeKey) bool {
	flag := keyToHasFlag(key)
	if flag == flagEmpty {
		return false
	}
	return node.flags.has(flag)
}

func (node *fullNode) IsNullFor(key NodeKey) bool {
	flag := keyToIsNullFlag(key)
	if flag == flagEmpty {
		return false
	}
	return node.flags.has(flag)
}

func (node *fullNode) NullableFor(key NodeKey) bool {
	flag := keyToNullableFlag(key)
	if flag == flagEmpty {
		return false
	}
	return node.flags.has(flag)
}

func (node *fullNode) ClearWhenEnterFor(key NodeKey) bool {
	flag := keyToClearWhenEnterFlag(key)
	if flag == flagEmpty {
		return false
	}
	return node.flags.has(flag)
}

func (node *fullNode) ModifyTime() ModifyTime {
	return node.modifyTime
}

func (node *fullNode) Desc() string {
	return node.descValue
}

func (node *fullNode) Int() int64 {
	return node.intValue
}

func (node *fullNode) Float() float64 {
	return node.floatValue
}

func (node *fullNode) Bool() bool {
	return node.boolValue
}

func (node *fullNode) String() string {
	return node.stringValue
}

func (node *fullNode) Obj() NodeObj {
	return node.objValue
}

func (node *fullNode) List() NodeList {
	return node.listValue
}

func (node *fullNode) ObjPrototype() *Node {
	return node.objPrototype
}

func (node *fullNode) ListPrototype() *Node {
	return node.listPrototype
}

// <<----- readonly methods end ----->>

// <<<==== side-effect methods begin ====>>>

func (node *fullNode) SetNullFor(key NodeKey, value bool) InnerNode {
	flag := keyToIsNullFlag(key)
	if value {
		node.flags.add(flag)
	} else {
		node.flags.delete(flag)
	}
	return node
}

func (node *fullNode) SetNullableFor(key NodeKey, value bool) InnerNode {
	flag := keyToNullableFlag(key)
	if value {
		node.flags.add(flag)
	} else {
		node.flags.delete(flag)
	}
	return node
}

func (node *fullNode) SetClearWhenEnterFor(key NodeKey, value bool) InnerNode {
	flag := keyToClearWhenEnterFlag(key)
	if flag == flagEmpty {
		return node
	}
	if value {
		node.flags.add(flag)
	} else {
		node.flags.delete(flag)
	}
	return node
}

func (node *fullNode) Delete(key NodeKey) InnerNode {
	flag := keyToHasFlag(key)
	if flag == flagEmpty {
		return node
	}
	node.flags.delete(flag)
	return node
}

func (node *fullNode) SetModifyTime(time ModifyTime) InnerNode {
	node.modifyTime = time
	return node
}

func (node *fullNode) SetDesc(value string) InnerNode {
	node.descValue = value
	node.flags.add(flagHasDesc)
	return node
}

func (node *fullNode) SetInt(value int64) InnerNode {
	node.intValue = value
	node.flags.add(flagHasInt)
	return node
}

func (node *fullNode) SetFloat(value float64) InnerNode {
	node.floatValue = value
	node.flags.add(flagHasFloat)
	return node
}

func (node *fullNode) SetBool(value bool) InnerNode {
	node.boolValue = value
	node.flags.add(flagHasBool)
	return node
}

func (node *fullNode) SetString(value string) InnerNode {
	node.stringValue = value
	node.flags.add(flagHasString)
	return node
}

func (node *fullNode) SetObj(value NodeObj) InnerNode {
	node.objValue = value
	node.flags.add(flagHasObj)
	return node
}

func (node *fullNode) SetList(value NodeList) InnerNode {
	node.listValue = value
	node.flags.add(flagHasList)
	return node
}

func (node *fullNode) SetObjPrototype(value *Node) InnerNode {
	node.objPrototype = value
	node.flags.add(flagHasObjPrototype)
	return node
}

func (node *fullNode) SetListPrototype(value *Node) InnerNode {
	node.listPrototype = value
	node.flags.add(flagHasListPrototype)
	return node
}

// <<----- side-effect methods begin ----->>
