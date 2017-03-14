// Package hash contains types for hashing and consistent hash rings
package hash

import (
	"fmt"

	farm "github.com/dgryski/go-farm"
	"github.com/tmrts/teracache/node"
	"github.com/tmrts/teracache/router/hash/tree"
)

const separator = '~'

type Func func(b []byte) int

type Ring interface {
	Insert(node.Type)
	Remove(node.Type)

	Search(string) node.Type
	SearchN(string, int) []node.Type
}

type ring struct {
	hash Func

	// virtualNodeCount is set to the number of extra nodes representing a
	// single unique entry. If a node is added to ring with 4 virtual nodes,
	// there will be 5 different keys that map to the same entry.
	virtualNodeCount int

	store tree.Set
}

// NewRing creates a new Ring with the requested amount of virtual nodes that
// are used to balance the hash ring.
func NewRing(vN int) Ring {
	return &ring{
		virtualNodeCount: vN,
		store:            tree.NewThreadSafeSet(),
		hash:             func(b []byte) int { return int(farm.Hash32(b)) },
	}
}

func (r *ring) Insert(node node.Type) {
	baseHash := r.hash(node.Bytes())

	// Normally a hash ring is susceptible to clustering due to the input
	// distribution, however we can sidestep this problem by creating a
	// good, empirically-tested, amount of virtual nodes to achieve an
	// evenly distributed hash ring.
	for i := 0; i < 1+r.virtualNodeCount; i++ {
		hash := baseHash ^ r.hash([]byte(fmt.Sprint(separator, i)))

		// Flyweight pattern is used here by inserting the reference to
		// the object instead of the object itself.
		r.store.Insert(hash, &node)
	}
}

func (r *ring) Remove(node node.Type) {
	baseHash := r.hash(node.Bytes())

	for i := 0; i < 1+r.virtualNodeCount; i++ {
		hash := baseHash ^ r.hash([]byte(fmt.Sprint(separator, i)))

		r.store.Remove(hash)
	}
}

// Search returns the node that is responsible for the given key.
func (r *ring) Search(key string) node.Type {
	return r.SearchN(key, 1)[0]
}

// SearchN returns at most N nodes that are responsible for the given key.
// It performs a best-effort search that returns at least 1, at most N keys.
func (r *ring) SearchN(key string, n int) []node.Type {
	hash := r.hash([]byte(key))

	return r.store.SearchN(hash, n)
}
