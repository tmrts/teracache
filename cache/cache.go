package cache

import (
	"sync"

	"github.com/hashicorp/golang-lru/simplelru"
	"github.com/tmrts/hordecache/payload"
	"github.com/tmrts/hordecache/with"
)

type Interface interface {
	Get(string) (payload.Payload, bool)

	Add(string, payload.Payload)

	Purge()
}

// EvictionNotice is invoked when an entry is evicted from the cache.
type EvictionNotice func(string, payload.Payload)

type lru struct {
	mu *sync.RWMutex

	cache *simplelru.LRU
}

// NewLRU creates a thread-safe cache with LRU as its eviction policy.
// Entries are evicted when the size is about to exceed the capacity.
// EvictionNotice callback is used whenever an entry is being evicted.
func NewLRU(capacity int, notify EvictionNotice) Interface {
	if capacity < 0 {
		panic("cache capacity must be non-negative!")
	}

	genericNotice := func(k, v interface{}) {
		notify(k.(string), v.(payload.Payload))
	}

	cache, _ := simplelru.NewLRU(capacity, genericNotice)

	return &lru{
		mu:    new(sync.RWMutex),
		cache: cache,
	}
}

func (c *lru) Add(k string, v payload.Payload) {
	_ = with.ReadLock(c.mu, func() error {
		_ = c.cache.Add(k, v)

		return nil
	})
}

func (c *lru) Get(k string) (payload.Payload, bool) {
	var (
		p  interface{}
		ok bool
	)

	// FIXME(tmrts): eviction callback is invoked inside the lock.
	_ = with.ReadLock(c.mu, func() error {
		p, ok = c.cache.Get(k)

		return nil
	})

	if !ok {
		return nil, false
	}

	return p.(payload.Payload), true
}

func (c *lru) Purge() {
	// TODO(tmrts): use an atomic reference CAS
	_ = with.Lock(c.mu, func() error {
		c.cache.Purge()

		return nil
	})
}
