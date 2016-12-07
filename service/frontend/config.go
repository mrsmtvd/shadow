package frontend

import (
	"github.com/kihamo/shadow/resource/config"
)

func (s *FrontendService) GetConfigVariables() []config.Variable {
	return []config.Variable{
		{
			Key:   "frontend.host",
			Value: "localhost",
			Usage: "Frontend host",
		},
		{
			Key:   "frontend.port",
			Value: 8080,
			Usage: "Frontend port number",
		},
		{
			Key:   "frontend.auth-user",
			Value: "",
			Usage: "User login",
		},
		{
			Key:   "frontend.auth-password",
			Value: "",
			Usage: "User password",
		},
	}
}
