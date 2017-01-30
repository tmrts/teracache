package service

import (
	"context"
	"fmt"
	"math/rand"
	"net"

	oldcontext "golang.org/x/net/context"

	"github.com/tmrts/hordecache/payload"
	"google.golang.org/grpc"
)

type cache interface {
	Get(context.Context, string) (payload.Payload, error)
}

type cacheServer struct {
	// TODO(tmrts): "cache" is written all over, find better names!

	cache cache
}

func NewServer(port int, c cache) (CacheServer, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	cs := &cacheServer{c}

	grpcServer := grpc.NewServer()

	RegisterCacheServer(grpcServer, cs)

	go grpcServer.Serve(lis)

	return cs, nil
}

func (c *cacheServer) Get(ctx oldcontext.Context, e *Entry) (*Payload, error) {
	// Convert ctx to new

	p, err := c.cache.Get(ctx, e.Key)
	if err != nil {
		return nil, err
	}

	return &Payload{
		ShouldCache: rand.Intn(10) == 0,
		Blob:        []byte(p),
	}, nil
}
