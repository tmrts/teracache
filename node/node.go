package node

import (
	"fmt"
	"net"
)

type Address struct {
	IP   net.IP
	Port uint16
}

func (addr *Address) String() string {
	return fmt.Sprint(addr.IP, addr.Port)
}

type Type struct {
	Identity string
	Address  Address
}

func (n *Type) Bytes() []byte {
	return []byte(fmt.Sprintf("%s-%s", n.Identity, n.Address))
}
