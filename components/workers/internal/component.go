package internal

import (
	"context"
	"sync"

	ws "github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/dispatcher"
	"github.com/kihamo/go-workers/listener"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/i18n"
	"github.com/kihamo/shadow/components/logging"
	"github.com/kihamo/shadow/components/metrics"
	"github.com/kihamo/shadow/components/workers"
)

type Component struct {
	application shadow.Application
	logger      logging.Logger

	mutex              sync.RWMutex
	dispatcher         *dispatcher.SimpleDispatcher
	lockedListenersIds []string
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
	c.lockedListenersIds = []string{}

	return nil
}

func (c *Component) Run(a shadow.Application, ready chan<- struct{}) error {
	c.logger = logging.DefaultLazyLogger(c.Name())

	<-a.ReadyComponent(config.ComponentName)
	cfg := a.GetComponent(config.ComponentName).(config.Component)

	c.dispatcher.SetTickerExecuteTasksDuration(cfg.Duration(workers.ConfigTickerExecuteTasksDuration))

	l := listener.NewFunctionListener(c.listenerLogging)
	l.SetName(c.Name() + ".logging")
	c.AddLockedListener(l.Id())

	c.AddListenerByEvents([]ws.Event{ws.EventAll}, l)

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

func (c *Component) GetLockedListeners() []string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	tmp := make([]string, len(c.lockedListenersIds))
	copy(tmp, c.lockedListenersIds)

	return tmp
}

func (c *Component) AddLockedListener(id string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.lockedListenersIds = append(c.lockedListenersIds, id)
}
