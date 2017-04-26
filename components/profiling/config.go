package profiling

import (
	"github.com/kihamo/shadow/components/config"
)

const (
	ConfigDumpDirectory = ComponentName + ".dump_directory"
	ConfigGCPercent     = ComponentName + ".gc_percent"
)

func (c *Component) GetConfigVariables() []config.Variable {
	return []config.Variable{
		{
			Key:      ConfigDumpDirectory,
			Usage:    "Path to trace dumps directory",
			Default:  "./",
			Type:     config.ValueTypeString,
			Editable: true,
		},
		{
			Key:      ConfigGCPercent,
			Usage:    "Sets the garbage collection target percentage",
			Default:  100,
			Type:     config.ValueTypeInt,
			Editable: true,
		},
	}
}

func (c *Component) GetConfigWatchers() map[string][]config.Watcher {
	return map[string][]config.Watcher{
		ConfigGCPercent: {c.watchGCPercent},
	}
}

func (c *Component) watchGCPercent(_ string, _ interface{}, _ interface{}) {
	c.initGCPercent()
}
