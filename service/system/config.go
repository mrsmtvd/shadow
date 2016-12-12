package system

import (
	"github.com/kihamo/shadow/resource/config"
)

const (
	ConfigSystemTimezone = "system.timezone"
)

func (s *SystemService) GetConfigVariables() []config.Variable {
	return []config.Variable{
		{
			Key:     ConfigSystemTimezone,
			Default: "Local",
			Usage:   "System timezone",
			Type:    config.ValueTypeString,
		},
	}
}
