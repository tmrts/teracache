package node

import "fmt"

type Type struct {
	Identity string
	Address  string
}

func (n *Type) Bytes() []byte {
	return []byte(fmt.Sprintf("%s-%s", n.Identity, n.Address))
}
