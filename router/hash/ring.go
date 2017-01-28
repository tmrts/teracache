package hash

import (
	"fmt"

	farm "github.com/dgryski/go-farm"
	"github.com/tmrts/hordecache/node"
	"github.com/tmrts/hordecache/router/hash/tree"
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
	hash             Func
	virtualNodeCount int

	store tree.Set
}

func NewRing(vN int) Ring {
	return &ring{
		virtualNodeCount: vN,
		store:            tree.NewThreadSafeSet(),
		hash:             func(b []byte) int { return int(farm.Hash32(b)) },
	}
}

func (r *ring) Insert(node node.Type) {
	baseHash := r.hash(node.Bytes())

	for i := 0; i < r.virtualNodeCount; i++ {
		hash := baseHash ^ r.hash([]byte(fmt.Sprint(separator, i)))

		r.store.Insert(hash, &node)
	}
}

func (r *ring) Remove(node node.Type) {
	baseHash := r.hash(node.Bytes())

	for i := 0; i < r.virtualNodeCount; i++ {
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
