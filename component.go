package shadow

import (
	"sync"
	"sync/atomic"
)

type Component interface {
	Name() string
	Version() string
	Run(Application, chan<- struct{}) error
}

type Dependency struct {
	Name     string
	Required bool
}

type ComponentInit interface {
	Init(Application) error
}

type ComponentDependency interface {
	Dependencies() []Dependency
}

type ComponentShutdown interface {
	Shutdown() error
}

type component struct {
	sync.RWMutex
	order    int64
	instance Component

	ready    int64
	done     int64
	error    atomic.Value
	watchers []chan struct{}
}

func newComponent(instance Component) *component {
	return &component{
		instance: instance,
		watchers: make([]chan struct{}, 0),
	}
}

func (c *component) Order() int64 {
	return atomic.LoadInt64(&c.order)
}

func (c *component) SetOrder(value int64) {
	atomic.StoreInt64(&c.order, value)
}

func (c *component) Ready() bool {
	return atomic.LoadInt64(&c.ready) == 1
}

func (c *component) Done() bool {
	return atomic.LoadInt64(&c.done) == 1
}

func (c *component) Error() error {
	if v := c.error.Load(); v != nil {
		return v.(error)
	}

	return nil
}

func (c *component) Run(a Application, result chan<- *component) {
	chReady := make(chan struct{}, 1)
	chDone := make(chan struct{}, 1)
	chError := make(chan error, 1)

	defer func() {
		close(chReady)
		close(chDone)
		close(chError)
	}()

	go func() {
		if err := c.instance.Run(a, chReady); err != nil {
			chError <- err
		} else {
			chDone <- struct{}{}
		}
	}()

	for {
		select {
		case <-chReady:
			atomic.StoreInt64(&c.ready, 1)

			c.notify()

		case err := <-chError:
			atomic.StoreInt64(&c.ready, 0)
			atomic.StoreInt64(&c.done, 1)
			c.error.Store(err)

			result <- c
			c.notify()

			return

		case <-chDone:
			atomic.StoreInt64(&c.ready, 1)
			atomic.StoreInt64(&c.done, 1)

			result <- c
			c.notify()

			return
		}
	}
}

func (c *component) notify() {
	c.Lock()
	tmp := append([]chan struct{}(nil), c.watchers...)
	c.watchers = make([]chan struct{}, 0)
	c.Unlock()

	go func(w []chan struct{}) {
		for _, ch := range w {
			close(ch)
		}
	}(tmp)
}

func (c *component) Watch() <-chan struct{} {
	ch := make(chan struct{}, 1)

	if c.Ready() {
		var wg sync.WaitGroup
		wg.Add(1)
		defer wg.Done()

		go func() {
			wg.Wait()
			close(ch)
		}()

		return ch
	} else {
		c.Lock()
		c.watchers = append(c.watchers, ch)
		c.Unlock()
	}

	return ch
}
