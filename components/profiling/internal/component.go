package internal

import (
	"runtime/debug"

	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/profiling"
	"github.com/kihamo/shadow/components/profiling/trace"
)

type Component struct {
	config config.Component
	routes []dashboard.Route
}

func (c *Component) GetName() string {
	return profiling.ComponentName
}

func (c *Component) GetVersion() string {
	return profiling.ComponentVersion
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
	c.config = a.GetComponent(config.ComponentName).(config.Component)

	return nil
}

func (c *Component) Run() error {
	c.initGCPercent()
	trace.LoadDumps(c.config.GetString(profiling.ConfigDumpDirectory))

	return nil
}

func (c *Component) initGCPercent() {
	debug.SetGCPercent(c.config.GetInt(profiling.ConfigGCPercent))
}
