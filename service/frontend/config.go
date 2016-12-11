package frontend

import (
	"github.com/kihamo/shadow/resource/config"
)

func (s *FrontendService) GetConfigVariables() []config.Variable {
	return []config.Variable{
		{
			Key:     "frontend.host",
			Default: "localhost",
			Usage:   "Frontend host",
			Type:    config.ValueTypeString,
		},
		{
			Key:     "frontend.port",
			Default: 8080,
			Usage:   "Frontend port number",
			Type:    config.ValueTypeInt,
		},
		{
			Key:   "frontend.auth-user",
			Usage: "User login",
			Type:  config.ValueTypeString,
		},
		{
			Key:   "frontend.auth-password",
			Usage: "User password",
			Type:  config.ValueTypeString,
		},
	}
}
