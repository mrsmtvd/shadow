package frontend

import (
	"github.com/kihamo/shadow/resource"
)

func (s *FrontendService) GetConfigVariables() []resource.ConfigVariable {
	return []resource.ConfigVariable{
		resource.ConfigVariable{
			Key:   "host",
			Value: "localhost",
			Usage: "Frontend host",
		},
		resource.ConfigVariable{
			Key:   "port",
			Value: int64(80),
			Usage: "Frontend port number",
		},
		resource.ConfigVariable{
			Key:   "auth-user",
			Value: "admin",
			Usage: "User login",
		},
		resource.ConfigVariable{
			Key:   "auth-password",
			Value: "password",
			Usage: "User password",
		},
	}
}
