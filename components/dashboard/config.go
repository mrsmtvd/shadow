package dashboard

import (
	"github.com/kihamo/shadow/components/config"
)

const (
	ConfigDashboardHost         = "dashboard.host"
	ConfigDashboardPort         = "dashboard.port"
	ConfigDashboardAuthUser     = "dashboard.auth-user"
	ConfigDashboardAuthPassword = "dashboard.auth-password"
)

func (c *Component) GetConfigVariables() []config.Variable {
	return []config.Variable{
		{
			Key:     ConfigDashboardHost,
			Default: "localhost",
			Usage:   "Frontend host",
			Type:    config.ValueTypeString,
		},
		{
			Key:     ConfigDashboardPort,
			Default: 8080,
			Usage:   "Frontend port number",
			Type:    config.ValueTypeInt,
		},
		{
			Key:      ConfigDashboardAuthUser,
			Usage:    "User login",
			Type:     config.ValueTypeString,
			Editable: true,
		},
		{
			Key:      ConfigDashboardAuthPassword,
			Usage:    "User password",
			Type:     config.ValueTypeString,
			Editable: true,
		},
	}
}
