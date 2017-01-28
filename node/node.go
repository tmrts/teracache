package node

import "fmt"

type Node struct {
	Identity string
	Address  string
}

func (n *Node) Bytes() []byte {
	return []byte(fmt.Sprintf("%s:%s", n.Identity, n.Address))
}
