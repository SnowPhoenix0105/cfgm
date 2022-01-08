package tree

import "testing"

func TestNode_ImplementReadWriteNode(t *testing.T) {
	var _ NodeReadWriter = &Node{}
}
