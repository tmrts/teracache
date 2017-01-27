package hordecache_test

import (
	"context"
	"fmt"
	"io"
	"testing"
)

func TestCreatesHordeCache(t *testing.T) {
	Store := map[string][]byte{
		"red":   []byte("#FF0000"),
		"green": []byte("#00FF00"),
		"blue":  []byte("#0000FF"),
	}

	backend := func(ctx context.Context, key string) (hordecache.Payloader, bool, error) {
		v, ok := Store[key]
		if !ok {
			return nil, false, nil
		}

		return hordecache.NewPayloader(v), true, nil
	}

	colors := hordecache.New("color-translation", 128, backend)

	ctx, cancel := context.WithCancel(nil)

	p, ok, err := colors.Get(ctx, "blue")
	if err != nil {
		panic(err)
	}
	cancel()

	if !ok {
		fmt.Println("value not available")
		return
	}

	var w io.Writer
	w.Write(p)
}
