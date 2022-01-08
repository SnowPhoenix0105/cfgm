package tree

type fullNodeFlag uint16

const (
	flagEmpty fullNodeFlag = 0
	flagDesc  fullNodeFlag = 1 << (iota - 1)
	flagInt
	flagFloat
	flagBool
	flagString
	flagObj
	flagList
	flagObjPrototype
	flagListPrototype
	// flagOverflow
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

func keyToFlag(key NodeKey) fullNodeFlag {
	switch key {
	case NodeKeyDesc:
		return flagDesc
	case NodeKeyInt:
		return flagInt
	case NodeKeyFloat:
		return flagFloat
	case NodeKeyBool:
		return flagBool
	case NodeKeyString:
		return flagString
	case NodeKeyObj:
		return flagObj
	case NodeKeyList:
		return flagList
	case NodeKeyObjPrototype:
		return flagObjPrototype
	case NodeKeyListPrototype:
		return flagListPrototype
	default:
		return flagEmpty
	}
}

type fullNode struct {
	validFlags    fullNodeFlag
	nullFlags     fullNodeFlag
	nullableFlags fullNodeFlag
	modifyTime    ModifyTime

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
	flag := keyToFlag(key)
	if flag == flagEmpty {
		return false
	}
	return node.validFlags.has(flag)
}

func (node *fullNode) IsNullFor(key NodeKey) bool {
	flag := keyToFlag(key)
	return node.nullFlags.has(flag)
}

func (node *fullNode) NullableFor(key NodeKey) bool {
	flag := keyToFlag(key)
	return node.nullableFlags.has(flag)
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
	flag := keyToFlag(key)
	if value {
		node.nullFlags.add(flag)
	} else {
		node.nullFlags.delete(flag)
	}
	return node
}

func (node *fullNode) SetNullableFor(key NodeKey, value bool) InnerNode {
	flag := keyToFlag(key)
	if value {
		node.nullableFlags.add(flag)
	} else {
		node.nullableFlags.delete(flag)
	}
	return node
}

func (node *fullNode) Delete(key NodeKey) InnerNode {
	flag := keyToFlag(key)
	if flag == flagEmpty {
		return node
	}
	node.validFlags.delete(flag)
	return node
}

func (node *fullNode) SetDesc(value string) InnerNode {
	node.descValue = value
	node.validFlags.add(flagDesc)
	return node
}

func (node *fullNode) SetModifyTime(time ModifyTime) InnerNode {
	node.modifyTime = time
	return node
}

func (node *fullNode) SetInt(value int64) InnerNode {
	node.intValue = value
	node.validFlags.add(flagInt)
	return node
}

func (node *fullNode) SetFloat(value float64) InnerNode {
	node.floatValue = value
	node.validFlags.add(flagFloat)
	return node
}

func (node *fullNode) SetBool(value bool) InnerNode {
	node.boolValue = value
	node.validFlags.add(flagBool)
	return node
}

func (node *fullNode) SetString(value string) InnerNode {
	node.stringValue = value
	node.validFlags.add(flagString)
	return node
}

func (node *fullNode) SetObj(value NodeObj) InnerNode {
	node.objValue = value
	node.validFlags.add(flagObj)
	return node
}

func (node *fullNode) SetList(value NodeList) InnerNode {
	node.listValue = value
	node.validFlags.add(flagList)
	return node
}

func (node *fullNode) SetObjPrototype(value *Node) InnerNode {
	node.objPrototype = value
	node.validFlags.add(flagObjPrototype)
	return node
}

func (node *fullNode) SetListPrototype(value *Node) InnerNode {
	node.listPrototype = value
	node.validFlags.add(flagListPrototype)
	return node
}

// <<----- side-effect methods begin ----->>
