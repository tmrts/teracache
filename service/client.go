package service

import (
	"context"

	"github.com/tmrts/hordecache/payload"
	"google.golang.org/grpc"
)

type Client func(context.Context) (payload.Payload, bool, error)

func NewClient(addr string, key string) Client {
	return func(ctx context.Context) (payload.Payload, bool, error) {
		// use context

		conn, err := grpc.Dial(addr, nil)
		if err != nil {
			return nil, false, err
		}
		defer conn.Close()

		client := NewCacheClient(conn)

		// TODO(tmrts): configure gRPC client using options
		resp, err := client.Get(ctx, &Entry{Key: key}, nil)
		if err != nil {
			return nil, false, err
		}

		return payload.Payload(resp.Blob), resp.ShouldCache, nil
	}
}
