package profiling

import (
	"github.com/kihamo/shadow/components/config"
)

const (
	ConfigProfilingDumpDirectory = ComponentName + ".dump_directory"
)

func (c *Component) GetConfigVariables() []config.Variable {
	return []config.Variable{
		{
			Key:      ConfigProfilingDumpDirectory,
			Usage:    "Path to trace dumps directory",
			Default:  "./",
			Type:     config.ValueTypeString,
			Editable: true,
		},
	}
}
