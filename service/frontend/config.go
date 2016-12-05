package frontend

import (
	"github.com/kihamo/shadow/resource/config"
)

func (s *FrontendService) GetConfigVariables() []config.ConfigVariable {
	return []config.ConfigVariable{
		config.ConfigVariable{
			Key:   "frontend.host",
			Value: "localhost",
			Usage: "Frontend host",
		},
		config.ConfigVariable{
			Key:   "frontend.port",
			Value: 8080,
			Usage: "Frontend port number",
		},
		config.ConfigVariable{
			Key:   "frontend.auth-user",
			Value: "",
			Usage: "User login",
		},
		config.ConfigVariable{
			Key:   "frontend.auth-password",
			Value: "",
			Usage: "User password",
		},
	}
}
