package profiling

import (
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/profiling/trace"
)

const (
	ComponentName = "profiling"
)

type Component struct {
	config *config.Component
}

func (c *Component) GetName() string {
	return ComponentName
}

func (c *Component) GetVersion() string {
	return "1.0.0"
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
	}
}

func (c *Component) Init(a shadow.Application) error {
	resourceConfig, err := a.GetComponent(config.ComponentName)
	if err != nil {
		return err
	}
	c.config = resourceConfig.(*config.Component)

	return nil
}

func (c *Component) Run() error {
	trace.LoadDumps(c.config.GetString(ConfigProfilingDumpDirectory))

	return nil
}
