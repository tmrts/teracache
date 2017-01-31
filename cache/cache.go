package cache

import (
	"github.com/hashicorp/golang-lru/simplelru"
	"github.com/tmrts/hordecache/payload"
)

type Key string

type Interface interface {
	Get(Key) (payload.Payload, bool)

	Add(Key, payload.Payload)

	Purge()
}

// EvictionNotice is invoked when an entry is evicted from the cache.
type EvictionNotice func(Key, payload.Payload)

type lru struct {
	cache *simplelru.LRU
}

func NewLRU(size int, notify EvictionNotice) Interface {
	if size < 0 {
		panic("cache size must be non-negative!")
	}

	genericNotice := func(k, v interface{}) {
		notify(k.(Key), v.(payload.Payload))
	}

	cache, _ := simplelru.NewLRU(size, genericNotice)

	return &lru{cache}
}

func (c *lru) Add(k Key, v payload.Payload) {
	_ = c.cache.Add(k, v)
}

func (c *lru) Get(k Key) (payload.Payload, bool) {
	p, ok := c.cache.Get(k)
	if !ok {
		return nil, false
	}

	return p.(payload.Payload), true
}

func (c *lru) Purge() {
	c.cache.Purge()
}
