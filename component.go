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

	status   int64
	runError atomic.Value

	watchers map[int64][]chan struct{}
}

func newComponent(instance Component) *component {
	return &component{
		instance: instance,
		watchers: make(map[int64][]chan struct{}, 0),
	}
}

func (c *component) Order() int64 {
	return atomic.LoadInt64(&c.order)
}

func (c *component) SetOrder(value int64) {
	atomic.StoreInt64(&c.order, value)
}

func (c *component) Status() componentStatus {
	return componentStatus(atomic.LoadInt64(&c.status))
}

func (c *component) setStatus(status componentStatus) {
	value := status.Int64()
	old := atomic.SwapInt64(&c.status, value)

	if old != value {
		c.notify(status)
	}
}

func (c *component) RunError() error {
	if v := c.runError.Load(); v != nil {
		return v.(error)
	}

	return nil
}

func (c *component) Run(a Application) {
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
		// компонент до завершения Run сообщил о своей готовности
		case <-chReady:
			if c.Status() != ComponentStatusShutdown {
				c.setStatus(ComponentStatusReady)
			}

			// не выходим, ждем следующий этап - полное завершение Run компонента

		// компонент не сообщал о готовности и Run вернул ошибку
		case err := <-chError:
			c.runError.Store(err)

			if c.Status() != ComponentStatusShutdown {
				c.setStatus(ComponentStatusRunFailed)
			}

			return

		// компонент не сообщал о готовности и Run успешно завершился
		case <-chDone:
			// в случае долгоиграющего компонента Run может разблокировать когда уже завершают приложение
			// в такой ситуации не надо менять уже установленный статус завершения
			if c.Status() != ComponentStatusShutdown {
				c.setStatus(ComponentStatusFinished)
			}

			return
		}
	}
}

func (c *component) Shutdown() (err error) {
	defer c.setStatus(ComponentStatusShutdown)

	closer, ok := c.instance.(ComponentShutdown)
	if !ok {
		return nil
	}

	return closer.Shutdown()
}

func (c *component) WatchStatus(status componentStatus) <-chan struct{} {
	ch := make(chan struct{}, 1)

	needClose := false
	current := c.Status()

	if status == current {
		needClose = true
	} else if current == ComponentStatusFinished {
		// если статус Finished то считается что компонент был в статусе Ready
		// если статус Finished то считается что компонент в статусе Shutdown
		needClose = status == ComponentStatusReady || status == ComponentStatusShutdown
	}

	if needClose {
		close(ch)
		return ch
	}

	// register watcher
	key := status.Int64()

	c.Lock()
	if _, ok := c.watchers[key]; !ok {
		c.watchers[key] = make([]chan struct{}, 0)
	}

	c.watchers[key] = append(c.watchers[key], ch)
	c.Unlock()

	return ch
}

func (c *component) notify(status componentStatus) {
	// если статус Running то это автоматически значит что достигнут статус Ready
	if status == ComponentStatusFinished {
		go func() {
			c.notify(ComponentStatusReady)
		}()
	}

	key := status.Int64()

	c.Lock()
	watchers, ok := c.watchers[key]

	// reset watchers
	if ok {
		delete(c.watchers, key)
	}
	c.Unlock()

	// watchers not found
	if !ok {
		return
	}

	go func(w []chan struct{}) {
		for _, ch := range w {
			close(ch)
		}
	}(watchers)
}
