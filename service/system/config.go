package system

import (
	"github.com/kihamo/shadow/resource/config"
)

func (s *SystemService) GetConfigVariables() []config.Variable {
	return []config.Variable{
		{
			Key:     "system.timezone",
			Default: "Local",
			Usage:   "System timezone",
			Type:    config.ValueTypeString,
		},
	}
}
