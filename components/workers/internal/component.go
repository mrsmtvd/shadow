package internal

import (
	"sync"

	ws "github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/dispatcher"
	"github.com/kihamo/go-workers/listener"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/i18n"
	"github.com/kihamo/shadow/components/logger"
	"github.com/kihamo/shadow/components/metrics"
	"github.com/kihamo/shadow/components/workers"
)

type Component struct {
	application shadow.Application
	config      config.Component
	logger      logger.Logger
	routes      []dashboard.Route

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
			Name: logger.ComponentName,
		},
		{
			Name: metrics.ComponentName,
		},
	}
}

func (c *Component) Init(a shadow.Application) error {
	c.application = a
	c.config = a.GetComponent(config.ComponentName).(config.Component)
	c.dispatcher = dispatcher.NewSimpleDispatcher()
	c.lockedListenersIds = []string{}

	return nil
}

func (c *Component) Run(wg *sync.WaitGroup) (err error) {
	c.logger = logger.NewOrNop(c.Name(), c.application)

	c.dispatcher.SetTickerExecuteTasksDuration(c.config.Duration(workers.ConfigTickerExecuteTasksDuration))

	l := listener.NewFunctionListener(c.listenerLogging)
	l.SetName(c.Name() + ".logging")
	c.AddLockedListener(l.Id())

	c.AddListenerByEvents([]ws.Event{ws.EventAll}, l)

	for i := 1; i <= c.config.Int(workers.ConfigWorkersCount); i++ {
		c.AddSimpleWorker()
	}

	go func() {
		defer wg.Done()
		c.dispatcher.Run()
	}()

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
