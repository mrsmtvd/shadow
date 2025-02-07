package internal

import (
	"runtime"

	"github.com/mrsmtvd/shadow/components/config"
	"github.com/mrsmtvd/shadow/components/profiling"
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
		config.NewVariable(profiling.ConfigProfileBlockRate, config.ValueTypeInt).
			WithUsage("Controls the fraction of goroutine blocking events that are reported in the blocking profile").
			WithEditable(true).
			WithDefault(10),
		config.NewVariable(profiling.ConfigProfileMutexFraction, config.ValueTypeInt).
			WithUsage("Controls the fraction of mutex contention events that are reported in the mutex profile").
			WithEditable(true).
			WithDefault(10),
	}
}

func (c *Component) ConfigWatchers() []config.Watcher {
	return []config.Watcher{
		config.NewWatcher([]string{profiling.ConfigGCPercent}, c.watchGCPercent),
		config.NewWatcher([]string{profiling.ConfigGoMaxProc}, c.watchGoMaxProc),
		config.NewWatcher([]string{profiling.ConfigProfileBlockRate}, c.watchProfileBlockRate),
		config.NewWatcher([]string{profiling.ConfigProfileMutexFraction}, c.watchProfileMutexFraction),
	}
}

func (c *Component) watchGCPercent(_ string, new interface{}, _ interface{}) {
	c.initGCPercent(new.(int))
}

func (c *Component) watchGoMaxProc(_ string, new interface{}, _ interface{}) {
	c.initGoMaxProc(new.(int))
}

func (c *Component) watchProfileBlockRate(_ string, new interface{}, _ interface{}) {
	c.initBlockProfile(new.(int))
}

func (c *Component) watchProfileMutexFraction(_ string, new interface{}, _ interface{}) {
	c.initMutexProfile(new.(int))
}
