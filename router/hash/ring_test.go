package hash_test

import (
	"testing"

	"github.com/tmrts/teracache/node"
	"github.com/tmrts/teracache/router/hash"
)

func TestSearchesRingWithoutVirtualNodes(t *testing.T) {
	ring := hash.NewRing(0)

	ring.Insert(node.Type{Identity: "test-node"})

	nod := ring.Search("some-key")
	if nod.Identity != "test-node" {
		t.Fatalf("ring.Search('some-key') expected 'test-node', got %q", nod.Identity)
	}
}
