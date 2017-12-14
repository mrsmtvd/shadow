package internal

import (
	"fmt"
	"sync"
	"time"

	ws "github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/dispatcher"
	"github.com/kihamo/go-workers/task"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/logger"
	"github.com/kihamo/shadow/components/metrics"
	"github.com/kihamo/shadow/components/workers"
)

type Component struct {
	application shadow.Application
	config      config.Component
	logger      logger.Logger
	routes      []dashboard.Route

	dispatcher ws.Dispatcher
}

func (c *Component) GetName() string {
	return workers.ComponentName
}

func (c *Component) GetVersion() string {
	return workers.ComponentVersion
}

func (c *Component) GetDependencies() []shadow.Dependency {
	return []shadow.Dependency{
		{
			Name:     config.ComponentName,
			Required: true,
		},
		{
			Name: dashboard.ComponentName,
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

	return nil
}

func (c *Component) Run(wg *sync.WaitGroup) (err error) {
	c.logger = logger.NewOrNop(c.GetName(), c.application)

	c.AddListener(ws.EventIdWorkerAdd, c.listenWorkerAdd)
	c.AddListener(ws.EventIdWorkerRemove, c.listenWorkerRemove)
	c.AddListener(ws.EventIdTaskAdd, c.listenTaskAdd)
	c.AddListener(ws.EventIdTaskRemove, c.listenTaskRemove)
	c.AddListener(ws.EventIdListenerAdd, c.listenListenerAdd)
	c.AddListener(ws.EventIdListenerRemove, c.listenListenerRemove)
	c.AddListener(ws.EventIdTaskExecuteStart, c.listenTaskExecuteStart)
	c.AddListener(ws.EventIdTaskExecuteStop, c.listenTaskExecuteStop)
	c.AddListener(ws.EventIdDispatcherStatusChanged, c.listenDispatcherStatusChanged)
	c.AddListener(ws.EventIdWorkerStatusChanged, c.listenWorkerStatusChanged)
	c.AddListener(ws.EventIdTaskStatusChanged, c.listenTaskStatusChanged)

	for i := 1; i <= c.config.GetInt(workers.ConfigWorkersCount); i++ {
		c.AddSimpleWorker()
	}

	go func() {
		defer wg.Done()
		c.dispatcher.Run()
	}()

	for i := 0; i < 10; i++ {
		t := task.NewFunctionTask(func() (interface{}, error) {
			time.Sleep(time.Hour * 10)

			return nil, fmt.Errorf("Error in time %s", time.Now().String())
		})
		t.SetName("test")
		t.SetPriority(int64(i))
		t.SetRepeats(-2)
		c.AddTask(t)
	}

	return nil
}
