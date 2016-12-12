package frontend

import (
	"github.com/kihamo/shadow/resource/config"
)

const (
	ConfigFrontendHost         = "frontend.host"
	ConfigFrontendPort         = "frontend.port"
	ConfigFrontendAuthUser     = "frontend.auth-user"
	ConfigFrontendAuthPassword = "frontend.auth-password"
)

func (s *FrontendService) GetConfigVariables() []config.Variable {
	return []config.Variable{
		{
			Key:     ConfigFrontendHost,
			Default: "localhost",
			Usage:   "Frontend host",
			Type:    config.ValueTypeString,
		},
		{
			Key:     ConfigFrontendPort,
			Default: 8080,
			Usage:   "Frontend port number",
			Type:    config.ValueTypeInt,
		},
		{
			Key:      ConfigFrontendAuthUser,
			Usage:    "User login",
			Type:     config.ValueTypeString,
			Editable: true,
		},
		{
			Key:      ConfigFrontendAuthPassword,
			Usage:    "User password",
			Type:     config.ValueTypeString,
			Editable: true,
		},
	}
}

func (s *FrontendService) GetConfigWatchers() map[string][]config.Watcher {
	return map[string][]config.Watcher{
		ConfigFrontendAuthUser:     {s.watchAuth},
		ConfigFrontendAuthPassword: {s.watchAuth},
	}
}

func (s *FrontendService) watchAuth(_ interface{}, _ interface{}) {
	s.generateAuthToken()
}
