package hordecache_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/tmrts/hordecache"
	"github.com/tmrts/hordecache/payload"
)

const (
	KB = 1 << 10
)

func TestCreatesTopic(t *testing.T) {
	store := map[string][]byte{
		"red":   []byte("#FF0000"),
		"green": []byte("#00FF00"),
		"blue":  []byte("#0000FF"),
	}

	// creates a 1 KB stand-alone cache for the provider
	colors, err := hordecache.New(hordecache.Topic{
		ID:       "colors",
		Capacity: 1 * KB,
		Peers:    []string{},
		Provider: func(_ context.Context, key string) (payload.Payload, error) {
			v, ok := store[key]
			if !ok {
				return nil, fmt.Errorf("key %#q not found", key)
			}

			return v, nil
		},
	})
	if err != nil {
		t.Fatalf("hordecache.New() got error %#v", err)
	}

	key := "blue"

	data, err := colors.Get(nil, key)
	if err != nil {
		t.Fatalf("hordecache.Get(%q) got error -> %q", key, err)
	}

	if want, got := string(store[key]), string(data); want != got {
		t.Fatalf("hordecache.Get(%q) expected %q, got  %q", key, want, got)
	}
}
