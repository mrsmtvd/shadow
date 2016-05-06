package frontend

import (
	"github.com/kihamo/shadow/resource"
)

func (s *FrontendService) GetConfigVariables() []resource.ConfigVariable {
	return []resource.ConfigVariable{
		resource.ConfigVariable{
			Key:   "frontend.host",
			Value: "localhost",
			Usage: "Frontend host",
		},
		resource.ConfigVariable{
			Key:   "frontend.port",
			Value: int64(8080),
			Usage: "Frontend port number",
		},
		resource.ConfigVariable{
			Key:   "frontend.auth-user",
			Value: "",
			Usage: "User login",
		},
		resource.ConfigVariable{
			Key:   "frontend.auth-password",
			Value: "",
			Usage: "User password",
		},
	}
}
