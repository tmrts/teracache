// Package with contains context managers for objects that require
// initialization before and clean-up after the execution such as files, locks,
// connections.
package with

import (
	"io"
	"sync"
)

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

func Closer(c io.Closer, fn func() error) (err error) {
	defer c.Close()

	return fn()
}
