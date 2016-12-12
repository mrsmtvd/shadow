package logger

import (
	"github.com/kihamo/shadow/resource/config"
)

const (
	ConfigLoggerLevel = "logger.level"
)

func (r *Resource) GetConfigVariables() []config.Variable {
	return []config.Variable{
		{
			Key:     ConfigLoggerLevel,
			Default: 5,
			Usage:   "Log level",
			Type:    config.ValueTypeInt,
		},
	}
}
