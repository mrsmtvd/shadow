package internal

import (
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/profiling"
)

func (c *Component) GetConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariableItem(
			profiling.ConfigDumpDirectory,
			config.ValueTypeString,
			"./",
			"Path to trace dumps directory",
			true,
			nil,
			nil),
		config.NewVariableItem(
			profiling.ConfigGCPercent,
			config.ValueTypeInt,
			100,
			"Sets the garbage collection target percentage",
			true,
			nil,
			nil),
	}
}

func (c *Component) GetConfigWatchers() map[string][]config.Watcher {
	return map[string][]config.Watcher{
		profiling.ConfigGCPercent: {c.watchGCPercent},
	}
}

func (c *Component) watchGCPercent(_ string, _ interface{}, _ interface{}) {
	c.initGCPercent()
}
