package internal

import (
	"context"
	"sync"

	ws "github.com/mrsmtvd/go-workers"
	"github.com/mrsmtvd/go-workers/dispatcher"
	"github.com/mrsmtvd/shadow"
	"github.com/mrsmtvd/shadow/components/config"
	"github.com/mrsmtvd/shadow/components/dashboard"
	"github.com/mrsmtvd/shadow/components/i18n"
	"github.com/mrsmtvd/shadow/components/logging"
	"github.com/mrsmtvd/shadow/components/metrics"
	"github.com/mrsmtvd/shadow/components/workers"
)

type Component struct {
	application shadow.Application
	logger      logging.Logger

	mutex           sync.RWMutex
	dispatcher      *dispatcher.SimpleDispatcher
	lockedListeners []ws.ListenerWithEvents
}

func (c *Component) Name() string {
	return workers.ComponentName
}

func (c *Component) Version() string {
	return workers.ComponentVersion
}

func (c *Component) Dependencies() []shadow.Dependency {
	return []shadow.Dependency{
		{
			Name:     config.ComponentName,
			Required: true,
		},
		{
			Name: dashboard.ComponentName,
		},
		{
			Name: i18n.ComponentName,
		},
		{
			Name: logging.ComponentName,
		},
		{
			Name: metrics.ComponentName,
		},
	}
}

func (c *Component) Init(a shadow.Application) error {
	c.application = a

	c.dispatcher = dispatcher.NewSimpleDispatcher()
	c.lockedListeners = make([]ws.ListenerWithEvents, 0)

	return nil
}

func (c *Component) Run(a shadow.Application, ready chan<- struct{}) error {
	c.logger = logging.DefaultLazyLogger(c.Name())

	<-a.ReadyComponent(config.ComponentName)
	cfg := a.GetComponent(config.ComponentName).(config.Component)

	c.dispatcher.SetTickerExecuteTasksDuration(cfg.Duration(workers.ConfigTickerExecuteTasksDuration))

	if cfg.Bool(workers.ConfigListenersLoggingEnabled) {
		if l := c.newLoggingListener(); l != nil {
			c.addLockedListener(l)
		}
	}

	for i := 1; i <= cfg.Int(workers.ConfigWorkersCount); i++ {
		c.AddSimpleWorker()
	}

	ready <- struct{}{}

	return c.dispatcher.Run()
}

func (c *Component) Shutdown() error {
	if c.dispatcher.Status() == ws.DispatcherStatusProcess {
		if err := c.dispatcher.Cancel(); err != context.Canceled {
			return err
		}
	}

	return nil
}

func (c *Component) LockedListeners() []ws.ListenerWithEvents {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	tmp := make([]ws.ListenerWithEvents, len(c.lockedListeners))
	copy(tmp, c.lockedListeners)

	return tmp
}

func (c *Component) addLockedListener(listener ws.ListenerWithEvents) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.AddListener(listener)
	c.lockedListeners = append(c.lockedListeners, listener)
}

func (c *Component) removeLockedListener(name string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for i := len(c.lockedListeners) - 1; i >= 0; i-- {
		if c.lockedListeners[i].Name() == name {
			c.RemoveListener(c.lockedListeners[i])
			c.lockedListeners = append(c.lockedListeners[:i], c.lockedListeners[i+1:]...)
		}
	}
}
