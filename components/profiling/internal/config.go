package internal

import (
	"runtime"

	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/profiling"
)

func (c *Component) ConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariable(
			profiling.ConfigDumpDirectory,
			config.ValueTypeString,
			"./",
			"Path to trace dumps directory",
			true,
			"Others",
			nil,
			nil),
		config.NewVariable(
			profiling.ConfigGCPercent,
			config.ValueTypeInt,
			100,
			"Sets the garbage collection target percentage",
			true,
			"Others",
			nil,
			nil),
		config.NewVariable(
			profiling.ConfigGoMaxProc,
			config.ValueTypeInt,
			runtime.GOMAXPROCS(-1),
			"Sets the maximum number of CPUs that can be executing simultaneously",
			true,
			"Others",
			nil,
			nil),
	}
}

func (c *Component) ConfigWatchers() []config.Watcher {
	return []config.Watcher{
		config.NewWatcher([]string{profiling.ConfigGCPercent}, c.watchGCPercent),
		config.NewWatcher([]string{profiling.ConfigGoMaxProc}, c.watchGoMaxProc),
	}
}

func (c *Component) watchGCPercent(_ string, _ interface{}, _ interface{}) {
	c.initGCPercent()
}

func (c *Component) watchGoMaxProc(_ string, _ interface{}, _ interface{}) {
	c.initGoMaxProc()
}
