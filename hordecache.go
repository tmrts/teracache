package hordecache

import (
	"context"

	"github.com/tmrts/hordecache/payload"
	"github.com/tmrts/hordecache/router"
	"github.com/tmrts/hordecache/service"

	lru "github.com/tmrts/hordecache/cache"
)

type Provider func(ctx context.Context, key string) (payload.Payload, error)

type Interface interface {
	Get(context.Context, string) (payload.Payload, error)
}

type horde struct {
	lru lru.Interface

	provider Provider

	router router.Interface

	svc service.CacheServer
}

const (
	RouterPort  = 20274
	ServicePort = 20275
)

func New(capacity int, hosts []string, p Provider) (Interface, error) {
	// TODO(tmrts): utilize the eviction callback in LRU
	lruCache := lru.NewLRU(capacity, nil)

	r, err := router.New(RouterPort)
	if err != nil {
		return nil, err
	}

	if err := r.Join(hosts); err != nil {
		return nil, err
	}

	c := &horde{
		lru:      lruCache,
		provider: p,
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

func (c *horde) Get(ctx context.Context, key string) (payload.Payload, error) {
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
