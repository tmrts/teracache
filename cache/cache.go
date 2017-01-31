package cache

import "github.com/tmrts/hordecache/payload"

type Key string

type Interface interface {
	Get(Key) (payload.Payload, bool)

	Add(Key, payload.Payload)

	Purge()
}
