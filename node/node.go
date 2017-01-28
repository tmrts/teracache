package node

import "fmt"

type Address string

func (addr *Address) String() string {
	return fmt.Sprint(*addr)
}

type Type struct {
	identity string
	address  Address
}

func (n *Type) Bytes() []byte {
	return []byte(fmt.Sprintf("%s-%s", n.identity, n.address))
}

func (n *Type) Address() Address {
	return n.address
}
