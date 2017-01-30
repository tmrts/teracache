package node

import (
	"fmt"
	"net"
)

type Addr struct {
	IP   net.IP
	Port uint16
}

func (a *Addr) String() string {
	return fmt.Sprint(a.IP, a.Port)
}

type Type struct {
	Identity string
	Addr     Addr
}

func (n *Type) Bytes() []byte {
	return []byte(fmt.Sprintf("%s-%s", n.Identity, n.Addr))
}

func New(name string, ip net.IP, port uint16) Type {
	return Type{
		Identity: name,
		Addr:     Addr{ip, port},
	}
}
