package property

type Record struct {
	root *node
}

type node struct {
	value string
	sub   map[string]*node
}

func newNode() *node {
	return &node{
		value: "",
		sub:   make(map[string]*node),
	}
}
