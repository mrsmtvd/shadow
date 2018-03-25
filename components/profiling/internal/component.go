package internal

import (
	"expvar"
	"runtime"
	"runtime/debug"

	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/i18n"
	"github.com/kihamo/shadow/components/profiling"
	"github.com/kihamo/shadow/components/profiling/trace"
)

type Component struct {
	config config.Component
	routes []dashboard.Route
}

func (c *Component) Name() string {
	return profiling.ComponentName
}

func (c *Component) Version() string {
	return profiling.ComponentVersion
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
	}
}

func (c *Component) Init(a shadow.Application) error {
	c.config = a.GetComponent(config.ComponentName).(config.Component)

	expvar.Publish(c.Name()+".runtime", expvar.Func(expvarRuntime))

	return nil
}

func (c *Component) Run() error {
	c.initGCPercent()
	c.initGoMaxProc()

	trace.LoadDumps(c.config.String(profiling.ConfigDumpDirectory))

	return nil
}

func (c *Component) initGCPercent() {
	debug.SetGCPercent(c.config.Int(profiling.ConfigGCPercent))
}

func (c *Component) initGoMaxProc() {
	runtime.GOMAXPROCS(c.config.Int(profiling.ConfigGoMaxProc))
}
