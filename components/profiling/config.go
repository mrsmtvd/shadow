package profiling

import (
	"github.com/kihamo/shadow/components/config"
)

const (
	ConfigDumpDirectory = ComponentName + ".dump_directory"
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
	}
}
