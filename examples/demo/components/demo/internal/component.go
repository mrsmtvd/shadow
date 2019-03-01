package internal

import (
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/database"
	"github.com/kihamo/shadow/examples/demo/components/demo"
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
	}
}

func (c *Component) Run(a shadow.Application, _ chan<- struct{}) error {
	<-a.ReadyComponent(database.ComponentName)

	return nil
}
