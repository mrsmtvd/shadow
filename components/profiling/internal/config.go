package internal

import (
	"runtime"

	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/profiling"
)

func (c *Component) ConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariable(profiling.ConfigDumpDirectory, config.ValueTypeString).
			WithUsage("Path to trace dumps directory").
			WithEditable(true).
			WithDefault("./"),
		config.NewVariable(profiling.ConfigGCPercent, config.ValueTypeInt).
			WithUsage("Sets the garbage collection target percentage").
			WithEditable(true).
			WithDefault(100),
		config.NewVariable(profiling.ConfigGoMaxProc, config.ValueTypeInt).
			WithUsage("Sets the maximum number of CPUs that can be executing simultaneously").
			WithEditable(true).
			WithDefault(runtime.GOMAXPROCS(-1)),
	}
}

func (c *Component) ConfigWatchers() []config.Watcher {
	return []config.Watcher{
		config.NewWatcher([]string{profiling.ConfigGCPercent}, c.watchGCPercent),
		config.NewWatcher([]string{profiling.ConfigGoMaxProc}, c.watchGoMaxProc),
	}
}

func (c *Component) watchGCPercent(_ string, new interface{}, _ interface{}) {
	c.initGCPercent(new.(int))
}

func (c *Component) watchGoMaxProc(_ string, new interface{}, _ interface{}) {
	c.initGoMaxProc(new.(int))
}
