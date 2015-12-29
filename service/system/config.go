package system

import (
	"github.com/kihamo/shadow/resource"
)

func (s *SystemService) GetConfigVariables() []resource.ConfigVariable {
	return []resource.ConfigVariable{
		resource.ConfigVariable{
			Key:   "system.timezone",
			Value: "Local",
			Usage: "System timezone",
		},
	}
}
