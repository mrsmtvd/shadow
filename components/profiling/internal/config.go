package internal

import (
	"fmt"
	"runtime"

	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/profiling"
)

func (c *Component) GetConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariable(
			profiling.ConfigDumpDirectory,
			config.ValueTypeString,
			"./",
			"Path to trace dumps directory",
			true,
			"",
			nil,
			nil),
		config.NewVariable(
			profiling.ConfigGCPercent,
			config.ValueTypeInt,
			100,
			"Sets the garbage collection target percentage",
			true,
			"",
			nil,
			nil),
		config.NewVariable(
			profiling.ConfigGoMaxProc,
			config.ValueTypeInt,
			runtime.GOMAXPROCS(-1),
			fmt.Sprintf("Sets the maximum number of CPUs that can be executing simultaneously. Attention number of cpu is %d", runtime.NumCPU()),
			true,
			"",
			nil,
			nil),
	}
}

func (c *Component) GetConfigWatchers() []config.Watcher {
	return []config.Watcher{
		config.NewWatcher(profiling.ComponentName, []string{profiling.ConfigGCPercent}, c.watchGCPercent),
		config.NewWatcher(profiling.ComponentName, []string{profiling.ConfigGoMaxProc}, c.watchGoMaxProc),
	}
}

func (c *Component) watchGCPercent(_ string, _ interface{}, _ interface{}) {
	c.initGCPercent()
}

func (c *Component) watchGoMaxProc(_ string, _ interface{}, _ interface{}) {
	c.initGoMaxProc()
}
