package teracache

import (
	"context"

	"github.com/tmrts/teracache/payload"
	"github.com/tmrts/teracache/router"
	"github.com/tmrts/teracache/service"

	lru "github.com/tmrts/teracache/cache"
)

type Provider func(ctx context.Context, key string) (payload.Payload, error)

type Interface interface {
	Get(context.Context, string) (payload.Payload, error)
}

type topic struct {
	lru lru.Interface

	provider Provider

	router router.Interface

	svc service.CacheServer
}

const (
	RouterPort  = 20274
	ServicePort = 20275
)

type Topic struct {
	ID       string
	Capacity int
	Peers    []string
	Provider Provider
}

// New creates a cache instance that participates in the given topic. Once a Get
// request to the cache fails, the cache uses the topic provider function to
// retrieve the missing element. The cache is bootstrapped using the given peer
// addresses.
func New(t Topic) (Interface, error) {
	// TODO(tmrts): utilize the eviction callback in LRU
	lruCache := lru.NewLRU(t.Capacity, nil)

	r, err := router.New(RouterPort)
	if err != nil {
		return nil, err
	}

	if err := r.Join(t.Peers); err != nil {
		return nil, err
	}

	c := &topic{
		lru:      lruCache,
		provider: t.Provider,
		router:   r,
	}

	// FIXME(tmrts): needs restructuring
	svc, err := service.NewServer(ServicePort, c)
	if err != nil {
		return nil, err
	}

	c.svc = svc

	return c, nil
}

func (c *topic) Get(ctx context.Context, key string) (payload.Payload, error) {
	obj, ok := c.lru.Get(key)
	if ok {
		return obj, nil
	}

	// FIXME(tmrts): leaky abstraction, refactor at once!
	owner, ownedByMe, err := c.router.Route(key)
	if !ownedByMe {
		clnt := service.NewClient(owner.Addr.String(), key)

		// TODO(tmrts): utilize the context and request-scoped information.
		p, shouldCache, err := clnt(context.TODO())
		if err != nil {
			return nil, err
		}

		if shouldCache {
			defer c.lru.Add(key, p)
		}

		return p, nil
	}

	p, err := c.provider(context.TODO(), key)
	if err != nil {
		return nil, err
	}
	defer c.lru.Add(key, p)

	return p, nil
}
