package tree

type ReadonlyWalker interface {
	NodeReader

	TryEnterObj(key string) bool
	TryEnterList(index int) bool
	Exit()
}

type Walker interface {
	NodeWriter
	ReadonlyWalker

	EnterObj(key string)
	EnterList(index int)
}

func ReadFrom(root *Node) ReadonlyWalker {
	return WriteFrom(root, -1)
}

func WriteFrom(root *Node, time ModifyTime) Walker {
	return &walker{
		currentNode: root,
		stack:       make(nodeStack, 0),
		time:        time,
	}
}

type nodeStack []*Node

func (stack *nodeStack) push(node *Node) {
	*stack = append(*stack, node)
}

func (stack *nodeStack) pop() *Node {
	length := len(*stack)
	ret := (*stack)[length-1]
	*stack = (*stack)[:length-1]
	return ret
}

type walker struct {
	currentNode *Node
	stack       nodeStack
	time        ModifyTime
}

// <<<==== readonly methods begin ====>>>

func (walker *walker) TryEnterObj(key string) bool {
	if !walker.currentNode.Has(NodeKeyObj) {
		return false
	}
	obj := walker.currentNode.Obj()
	next, ok := obj[key]
	if !ok {
		return false
	}
	walker.stack.push(walker.currentNode)
	walker.currentNode = next
	return true
}

func (walker *walker) TryEnterList(index int) bool {
	if !walker.currentNode.Has(NodeKeyList) {
		return false
	}
	list := walker.currentNode.List()
	length := len(list)
	if length <= index {
		return false
	}
	walker.stack.push(walker.currentNode)
	walker.currentNode = list[index]
	return true
}

func (walker *walker) Exit() {
	walker.currentNode = walker.stack.pop()
}

func (walker *walker) Has(key NodeKey) bool {
	return walker.currentNode.Has(key)
}

func (walker *walker) IsNullFor(key NodeKey) bool {
	return walker.currentNode.IsNullFor(key)
}

func (walker *walker) NullableFor(key NodeKey) bool {
	return walker.currentNode.NullableFor(key)
}

func (walker *walker) ModifyTime() ModifyTime {
	return walker.currentNode.ModifyTime()
}

func (walker *walker) Desc() string {
	return walker.currentNode.Desc()
}

func (walker *walker) Int() int64 {
	return walker.currentNode.Int()
}

func (walker *walker) Float() float64 {
	return walker.currentNode.Float()
}

func (walker *walker) Bool() bool {
	return walker.currentNode.Bool()
}

func (walker *walker) String() string {
	return walker.currentNode.String()
}

func (walker *walker) Obj() NodeObj {
	return walker.currentNode.Obj()
}

func (walker *walker) List() NodeList {
	return walker.currentNode.List()
}

func (walker *walker) ObjPrototype() *Node {
	return walker.currentNode.ObjPrototype()
}

func (walker *walker) ListPrototype() *Node {
	return walker.currentNode.ListPrototype()
}

// <<----- readonly methods end ----->>

// <<<==== side-effect methods begin ====>>>

func (walker *walker) setModifyTimeForParentNodes() {
	length := len(walker.stack)
	for i := length - 1; i >= 0; i-- {
		node := walker.stack[i]
		if node.ModifyTime() == walker.time {
			return
		}
		node.SetModifyTime(walker.time)
	}
}

func (walker *walker) EnterObj(key string) {
	var next *Node
	if !walker.currentNode.Has(NodeKeyObj) {
		next = NewNode()
		walker.currentNode.SetObj(map[string]*Node{key: next})
		walker.currentNode.SetModifyTime(walker.time)
		walker.setModifyTimeForParentNodes()
	} else {
		obj := walker.currentNode.Obj()
		var ok bool
		next, ok = obj[key]
		if !ok {
			next = NewNode()
			obj[key] = next
			walker.currentNode.SetObj(obj)
			walker.currentNode.SetModifyTime(walker.time)
			walker.setModifyTimeForParentNodes()
		}
	}
	walker.stack.push(walker.currentNode)
	walker.currentNode = next
}

func (walker *walker) EnterList(index int) {
	var next *Node
	if !walker.currentNode.Has(NodeKeyList) {
		list := make(NodeList, index+1)
		for i := 0; i <= index; i++ {
			next = NewNode()
			list[i] = next
		}
		list[index] = next
		walker.currentNode.SetList(list)
		walker.currentNode.SetModifyTime(walker.time)
		walker.setModifyTimeForParentNodes()
	} else {
		list := walker.currentNode.List()
		length := len(list)
		if index < length {
			next = list[index]
		} else {
			for i := length; i <= index; i++ {
				next = NewNode()
				list = append(list, next)
			}
			walker.currentNode.SetList(list)
			walker.currentNode.SetModifyTime(walker.time)
			walker.setModifyTimeForParentNodes()
		}
	}
	walker.stack.push(walker.currentNode)
	walker.currentNode = next
}

func (walker *walker) SetNullFor(key NodeKey, value bool) {
	walker.currentNode.SetNullFor(key, value)
	walker.currentNode.SetModifyTime(walker.time)
	walker.setModifyTimeForParentNodes()
}

func (walker *walker) SetNullableFor(key NodeKey, value bool) {
	walker.currentNode.SetNullableFor(key, value)
	walker.currentNode.SetModifyTime(walker.time)
	walker.setModifyTimeForParentNodes()
}

func (walker *walker) Delete(key NodeKey) {
	walker.currentNode.Delete(key)
	walker.currentNode.SetModifyTime(walker.time)
	walker.setModifyTimeForParentNodes()
}

func (walker *walker) SetDesc(value string) {
	walker.currentNode.SetDesc(value)
	walker.currentNode.SetModifyTime(walker.time)
	walker.setModifyTimeForParentNodes()
}

func (walker *walker) SetModifyTime(time ModifyTime) {
	walker.currentNode.SetModifyTime(time)
}

func (walker *walker) SetInt(value int64) {
	walker.currentNode.SetInt(value)
	walker.currentNode.SetModifyTime(walker.time)
	walker.setModifyTimeForParentNodes()
}

func (walker *walker) SetFloat(value float64) {
	walker.currentNode.SetFloat(value)
	walker.currentNode.SetModifyTime(walker.time)
	walker.setModifyTimeForParentNodes()
}

func (walker *walker) SetBool(value bool) {
	walker.currentNode.SetBool(value)
	walker.currentNode.SetModifyTime(walker.time)
	walker.setModifyTimeForParentNodes()
}

func (walker *walker) SetString(value string) {
	walker.currentNode.SetString(value)
	walker.currentNode.SetModifyTime(walker.time)
	walker.setModifyTimeForParentNodes()
}

func (walker *walker) SetObj(value NodeObj) {
	walker.currentNode.SetObj(value)
	walker.currentNode.SetModifyTime(walker.time)
	walker.setModifyTimeForParentNodes()
}

func (walker *walker) SetList(value NodeList) {
	walker.currentNode.SetList(value)
	walker.currentNode.SetModifyTime(walker.time)
	walker.setModifyTimeForParentNodes()
}

func (walker *walker) SetObjPrototype(value *Node) {
	walker.currentNode.SetObjPrototype(value)
	walker.currentNode.SetModifyTime(walker.time)
	walker.setModifyTimeForParentNodes()
}

func (walker *walker) SetListPrototype(value *Node) {
	walker.currentNode.SetListPrototype(value)
	walker.currentNode.SetModifyTime(walker.time)
	walker.setModifyTimeForParentNodes()
}

// <<----- side-effect methods begin ----->>
