package internal

import (
	"context"
	"time"

	"github.com/mrsmtvd/go-workers/task"
	"github.com/mrsmtvd/shadow"
	"github.com/mrsmtvd/shadow/components/database"
	"github.com/mrsmtvd/shadow/components/logging"
	"github.com/mrsmtvd/shadow/components/workers"
	"github.com/mrsmtvd/shadow/examples/demo/components/demo"
)

type Component struct {
}

func (c *Component) Name() string {
	return demo.ComponentName
}

func (c *Component) Version() string {
	return demo.ComponentVersion
}

func (c *Component) Dependencies() []shadow.Dependency {
	return []shadow.Dependency{
		{
			Name:     database.ComponentName,
			Required: true,
		},
		{
			Name: logging.ComponentName,
		},
		{
			Name:     workers.ComponentName,
			Required: true,
		},
	}
}

func (c *Component) Run(a shadow.Application, _ chan<- struct{}) error {
	<-a.ReadyComponent(database.ComponentName)

	logger := logging.DefaultLazyLogger(c.Name())

	t := task.NewFunctionTask(func(_ context.Context) (interface{}, error) {
		logger.Error("Hello world! It's demo application")
		return nil, nil
	})
	t.SetName("task-" + c.Name())
	t.SetRepeats(-1)
	t.SetRepeatInterval(time.Second * 10)

	<-a.ReadyComponent(workers.ComponentName)
	a.GetComponent(workers.ComponentName).(workers.Component).AddTask(t)

	return nil
}
