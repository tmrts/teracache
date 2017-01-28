package with

import "sync"

func Lock(l sync.Locker, fn func() error) (err error) {
	l.Lock()
	defer l.Unlock()

	return fn()
}

type ReadLocker interface {
	RLocker() sync.Locker
}

func ReadLock(l ReadLocker, fn func() error) (err error) {
	return Lock(l.RLocker(), fn)
}
