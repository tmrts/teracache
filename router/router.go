package router

import (
	"fmt"

	"github.com/hashicorp/memberlist"
	"github.com/tmrts/hordecache/node"
	"github.com/tmrts/hordecache/router/hash"
)

type Interface interface {
	Join([]string) error

	LocalNode() node.Type

	Route(string) (node.Type, bool, error)
}

type router struct {
	ring hash.Ring

	list *memberlist.Memberlist
}

func New(port int) (Interface, error) {
	ring := hash.NewRing(50)

	cfg := memberlist.DefaultWANConfig()
	cfg.BindPort = port
	cfg.Name = fmt.Sprint(cfg.BindAddr, ":", cfg.BindPort)

	cfg.Events = newEventDelegate(ring)

	list, err := memberlist.Create(cfg)
	if err != nil {
		return nil, err
	}

	return &router{ring, list}, nil
}

func (r *router) Join(nodes []string) error {
	if _, err := r.list.Join(nodes); err != nil {
		return err
	}

	return nil
}

func memberToNode(member *memberlist.Node) node.Type {
	return node.New(member.Name, member.Addr, member.Port)
}

func (r *router) LocalNode() node.Type {
	// FIXME(tmrts): leaky abstraction, restructure to fix it
	m := r.list.LocalNode()

	return memberToNode(m)
}

func (r *router) Route(key string) (node.Type, bool, error) {
	owner := r.ring.Search(key)

	return owner, r.LocalNode().Identity == owner.Identity, nil
}

type delegate struct {
	ring hash.Ring
}

func newEventDelegate(r hash.Ring) memberlist.EventDelegate {
	return &delegate{r}
}

func (d *delegate) NotifyJoin(m *memberlist.Node) {
	d.ring.Insert(memberToNode(m))
}

func (d *delegate) NotifyLeave(m *memberlist.Node) {
	d.ring.Remove(memberToNode(m))
}

func (d *delegate) NotifyUpdate(node *memberlist.Node) {
	// TODO(tmrts): utilize the node metadata update notifications.
	return
}
