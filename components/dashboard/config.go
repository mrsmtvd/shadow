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

func (c *Component) GetConfigWatchers() map[string][]config.Watcher {
	return map[string][]config.Watcher{
		config.WatcherForAll: {c.watchConfig},
	}
}

func (c *Component) watchConfig(key string, newValue interface{}, oldValue interface{}) {
	c.logger.Infof("Change value for %s with '%v' to '%v'", key, oldValue, newValue)
}
