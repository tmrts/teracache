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
// Entries are evicted when the size is about to exceed the limit.
// EvictionNotice callback is used whenever an entry is being evicted.
func NewLRU(size int, notify EvictionNotice) Interface {
	if size < 0 {
		panic("cache size must be non-negative!")
	}

	genericNotice := func(k, v interface{}) {
		notify(k.(string), v.(payload.Payload))
	}

	cache, _ := simplelru.NewLRU(size, genericNotice)

	l := new(lru)
	l.cache = cache

	return l
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
