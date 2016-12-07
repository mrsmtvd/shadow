package logger

import (
	"github.com/kihamo/shadow/resource/config"
)

func (r *Resource) GetConfigVariables() []config.Variable {
	return []config.Variable{
		{
			Key:   "logger.level",
			Value: 5,
			Usage: "Log level",
		},
	}
}
