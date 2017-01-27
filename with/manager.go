package with

import "sync"

func Lock(l sync.Locker, fn func() error) (err error) {
	l.Lock()
	defer l.Unlock()

	return fn()
}
