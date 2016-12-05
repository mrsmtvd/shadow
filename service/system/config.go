package system

import (
	"github.com/kihamo/shadow/resource/config"
)

func (s *SystemService) GetConfigVariables() []config.ConfigVariable {
	return []config.ConfigVariable{
		config.ConfigVariable{
			Key:   "system.timezone",
			Value: "Local",
			Usage: "System timezone",
		},
	}
}
