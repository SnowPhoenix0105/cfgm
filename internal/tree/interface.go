package tree

type ModifyTime int8

type NodeKey int

const (
	NodeKeyInvalid NodeKey = iota
	NodeKeyDesc
	NodeKeyInt
	NodeKeyFloat
	NodeKeyBool
	NodeKeyString
	NodeKeyObj
	NodeKeyList
	NodeKeyObjPrototype
	NodeKeyListPrototype
)

func (key NodeKey) String() string {
	switch key {
	case NodeKeyInt:
		return "NodeKeyInt"
	case NodeKeyFloat:
		return "NodeKeyFloat"
	case NodeKeyBool:
		return "NodeKeyBool"
	case NodeKeyString:
		return "NodeKeyString"
	case NodeKeyObj:
		return "NodeKeyObj"
	case NodeKeyList:
		return "NodeKeyList"
	default:
		return "NodeKeyInvalid"
	}
}

type NodeObj map[string]*Node

type NodeList []*Node

type ReadableInnerNode interface {
	Has(key NodeKey) bool
	IsNullFor(key NodeKey) bool
	NullableFor(key NodeKey) bool
	ModifyTime() ModifyTime

	Desc() string
	Int() int64
	Float() float64
	Bool() bool
	String() string
	Obj() NodeObj
	List() NodeList
	ObjPrototype() *Node
	ListPrototype() *Node
}

type InnerNode interface {
	ReadableInnerNode
	Delete(key NodeKey) InnerNode
	SetNullFor(key NodeKey, value bool) InnerNode
	SetNullableFor(key NodeKey, value bool) InnerNode
	SetModifyTime(time ModifyTime) InnerNode

	SetDesc(value string) InnerNode
	SetInt(value int64) InnerNode
	SetFloat(value float64) InnerNode
	SetBool(value bool) InnerNode
	SetString(value string) InnerNode
	SetObj(value NodeObj) InnerNode
	SetList(value NodeList) InnerNode
	SetObjPrototype(value *Node) InnerNode
	SetListPrototype(value *Node) InnerNode
}

type NodeReader interface {
	ReadableInnerNode
}

type NodeWriter interface {
	Delete(key NodeKey)
	SetNullFor(key NodeKey, value bool)
	SetNullableFor(key NodeKey, value bool)
	SetModifyTime(time ModifyTime)

	SetDesc(value string)
	SetInt(value int64)
	SetFloat(value float64)
	SetBool(value bool)
	SetString(value string)
	SetObj(value NodeObj)
	SetList(value NodeList)
	SetObjPrototype(value *Node)
	SetListPrototype(value *Node)
}

type NodeReadWriter interface {
	NodeReader
	NodeWriter
}

type Node struct {
	Raw InnerNode
}

func NewNode() *Node {
	return &Node{
		Raw: newFullNode(),
	}
}

// <<<==== readonly methods begin ====>>>

func (node *Node) Has(key NodeKey) bool {
	return node.Raw.Has(key)
}

func (node *Node) IsNullFor(key NodeKey) bool {
	return node.Raw.IsNullFor(key)
}

func (node *Node) NullableFor(key NodeKey) bool {
	return node.Raw.NullableFor(key)
}

func (node *Node) ModifyTime() ModifyTime {
	return node.Raw.ModifyTime()
}

func (node *Node) Desc() string {
	return node.Desc()
}

func (node *Node) Int() int64 {
	return node.Raw.Int()
}

func (node *Node) Float() float64 {
	return node.Raw.Float()
}

func (node *Node) Bool() bool {
	return node.Raw.Bool()
}

func (node *Node) String() string {
	return node.Raw.String()
}

func (node *Node) Obj() NodeObj {
	return node.Raw.Obj()
}

func (node *Node) List() NodeList {
	return node.Raw.List()
}

func (node *Node) ObjPrototype() *Node {
	return node.Raw.ObjPrototype()
}

func (node *Node) ListPrototype() *Node {
	return node.Raw.ListPrototype()
}

// <<----- readonly methods end ----->>

// <<<==== side-effect methods begin ====>>>

func (node *Node) SetNullFor(key NodeKey, value bool) {
	node.Raw = node.Raw.SetNullFor(key, value)
}

func (node *Node) SetModifyTime(time ModifyTime) {
	node.Raw = node.Raw.SetModifyTime(time)
}

func (node *Node) SetNullableFor(key NodeKey, value bool) {
	node.Raw = node.Raw.SetNullableFor(key, value)
}

func (node *Node) Delete(key NodeKey) {
	node.Raw = node.Raw.Delete(key)
}

func (node *Node) SetDesc(value string) {
	node.Raw = node.Raw.SetDesc(value)
}

func (node *Node) SetInt(value int64) {
	node.Raw = node.Raw.SetInt(value)
}

func (node *Node) SetFloat(value float64) {
	node.Raw = node.Raw.SetFloat(value)
}

func (node *Node) SetString(value string) {
	node.Raw = node.Raw.SetString(value)
}

func (node *Node) SetBool(value bool) {
	node.Raw = node.Raw.SetBool(value)
}

func (node *Node) SetObj(value NodeObj) {
	node.Raw = node.Raw.SetObj(value)
}

func (node *Node) SetList(value NodeList) {
	node.Raw = node.Raw.SetList(value)
}

func (node *Node) SetObjPrototype(value *Node) {
	node.Raw = node.Raw.SetObjPrototype(value)
}

func (node *Node) SetListPrototype(value *Node) {
	node.Raw = node.Raw.SetListPrototype(value)
}

// <<----- side-effect methods begin ----->>
