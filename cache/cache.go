package cache

import "github.com/tmrts/hordecache/payload"

type Interface interface {
	Get(string) (payload.Payload, bool)
	Add(string, payload.Payload)

	Purge()
}
