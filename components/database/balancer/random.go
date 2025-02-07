package balancer

import (
	"math/rand"
	"sync"
	"sync/atomic"

	"github.com/mrsmtvd/shadow/components/database"
)

type Random struct {
	mutex     sync.RWMutex
	executors []database.Executor
	length    uint64
}

func NewRandom() *Random {
	return &Random{}
}

func (b *Random) Get() database.Executor {
	l := atomic.LoadUint64(&b.length)
	if l == 0 {
		return nil
	}

	i := 0
	if l > 1 {
		i = rand.Int() % int(l)
	}

	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return b.executors[i]
}

func (b *Random) Set(executors []database.Executor) {
	b.mutex.Lock()
	b.executors = executors
	b.mutex.Unlock()

	atomic.StoreUint64(&b.length, uint64(len(executors)))
}
