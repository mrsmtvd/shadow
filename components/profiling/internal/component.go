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

func (c *Component) Run(a shadow.Application, _ chan<- struct{}) error {
	expvar.Publish(c.Name()+".runtime", expvar.Func(expvarRuntime))

	<-a.ReadyComponent(config.ComponentName)
	cfg := a.GetComponent(config.ComponentName).(config.Component)

	c.initGCPercent(cfg.Int(profiling.ConfigGCPercent))
	c.initGoMaxProc(cfg.Int(profiling.ConfigGoMaxProc))
	c.initBlockProfile(cfg.Int(profiling.ConfigProfileBlockRate))
	c.initMutexProfile(cfg.Int(profiling.ConfigProfileMutexFraction))

	return trace.LoadDumps(cfg.String(profiling.ConfigDumpDirectory))
}

func (c *Component) initGCPercent(value int) {
	debug.SetGCPercent(value)
}

func (c *Component) initGoMaxProc(value int) {
	runtime.GOMAXPROCS(value)
}

func (c *Component) initBlockProfile(rate int) {
	runtime.SetBlockProfileRate(rate)
}

func (c *Component) initMutexProfile(fraction int) {
	runtime.SetMutexProfileFraction(fraction)
}
