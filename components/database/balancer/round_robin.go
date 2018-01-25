package balancer

import (
	"sync"
	"sync/atomic"

	"github.com/kihamo/shadow/components/database"
)

type RoundRobin struct {
	mutex     sync.RWMutex
	executors []database.Executor
	index     uint64
	length    uint64
}

func NewRoundRobin() *RoundRobin {
	return &RoundRobin{}
}

func (b *RoundRobin) Get() database.Executor {
	l := atomic.LoadUint64(&b.length)
	if l == 0 {
		return nil
	}

	atomic.AddUint64(&b.index, 1)
	i := atomic.LoadUint64(&b.index)
	if i >= l {
		i = 0
	}
	atomic.StoreUint64(&b.index, i)

	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return b.executors[i]
}

func (b *RoundRobin) Set(executors []database.Executor) {
	b.mutex.Lock()
	b.executors = executors
	b.mutex.Unlock()

	atomic.StoreUint64(&b.index, 0)
	atomic.StoreUint64(&b.length, uint64(len(executors)))
}
