package tree

import (
	"sort"
	"sync"

	"github.com/tmrts/hordecache/node"
	"github.com/tmrts/hordecache/with"
)

type Set interface {
	Insert(int, node.Type)
	Remove(int)

	Search(int) node.Type
	SearchN(int, int) []node.Type
}

type record struct {
	Key  int
	Data *node.Type
}

type recordSorter struct {
	records []record
}

func (r *recordSorter) Len() int { return len(r.records) }

func (r *recordSorter) Less(i, j int) bool { return r.records[i].Key < r.records[j].Key }

func (r *recordSorter) Swap(i, j int) { r.records[i], r.records[j] = r.records[j], r.records[i] }

// TODO(tmrts): change this when it becomes a bottleneck with:
//              1. AVL tree
//              2. vEB tree
//              3. Search tree that uses local-only changes for minimizing lock overhead
//              4. A read-copy-update tree
type ThreadSafeSet struct {
	lock *sync.RWMutex

	records []record
}

func NewThreadSafeSet() Set {
	return &ThreadSafeSet{
		lock: new(sync.RWMutex),
	}
}

func (s *ThreadSafeSet) Insert(key int, data node.Type) {
	with.Lock(s.lock, func() error {
		s.records = append(s.records, record{key, &data})

		sort.Sort(&recordSorter{s.records})

		return nil
	})
}

func (s *ThreadSafeSet) Remove(key int) {
	// TODO(tmrts): handle collisions
	with.Lock(s.lock, func() error {
		i := sort.Search(len(s.records), func(i int) bool { return s.records[i].Key >= key })
		if i >= 0 && i < len(s.records) {
			if s.records[i].Key == key {
				s.records = append(s.records[:i], s.records[i+1:]...)
			}
		}

		return nil
	})
}

func (s *ThreadSafeSet) Search(key int) node.Type {
	// FIXME(tmrts): assumes that set is never empty
	return s.SearchN(key, 1)[0]
}

func (s *ThreadSafeSet) SearchN(key, _ int) []node.Type {
	ch := make(chan *node.Type, 1)

	go with.ReadLock(s.lock, func() error {
		// TODO(tmrts): use a bucket search to make this O(1)
		i := sort.Search(len(s.records), func(i int) bool { return s.records[i].Key >= key })

		if i >= 0 && i < len(s.records) {
			ch <- s.records[i].Data
		} else {
			// FIXME(tmrts): will lead to minimal request clustering
			ch <- s.records[0].Data
		}

		return nil
	})

	n := <-ch
	return []node.Type{*n}
}
