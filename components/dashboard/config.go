package dashboard

import (
	"github.com/kihamo/shadow/components/config"
)

const (
	ConfigFrontendHost         = "frontend.host"
	ConfigFrontendPort         = "frontend.port"
	ConfigFrontendAuthUser     = "frontend.auth-user"
	ConfigFrontendAuthPassword = "frontend.auth-password"
)

func (c *Component) GetConfigVariables() []config.Variable {
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

func (c *Component) GetConfigWatchers() map[string][]config.Watcher {
	return map[string][]config.Watcher{
		ConfigFrontendAuthUser:     {c.watchAuthUser},
		ConfigFrontendAuthPassword: {c.watchAuthPassword},
	}
}

func (c *Component) watchAuthUser(newValue interface{}, _ interface{}) {
	c.generateAuthToken(newValue.(string), c.config.GetString(ConfigFrontendAuthPassword))
}

func (c *Component) watchAuthPassword(newValue interface{}, _ interface{}) {
	c.generateAuthToken(c.config.GetString(ConfigFrontendAuthUser), newValue.(string))
}
