package internal

import (
	"github.com/kihamo/shadow/components/annotations"
	"github.com/kihamo/shadow/components/config"
)

func (c *Component) ConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariable(
			annotations.ConfigStorageGrafanaEnabled,
			config.ValueTypeBool,
			false,
			"Enabled",
			true,
			"Grafana storage",
			nil,
			nil),
		config.NewVariable(
			annotations.ConfigStorageGrafanaAddress,
			config.ValueTypeString,
			nil,
			"Grafana address of HTTP API in format http://host:port",
			true,
			"Grafana storage",
			nil,
			nil),
		config.NewVariable(
			annotations.ConfigStorageGrafanaApiKey,
			config.ValueTypeString,
			nil,
			"API key. No need if username and password is set",
			true,
			"Grafana storage",
			nil,
			nil),
		config.NewVariable(
			annotations.ConfigStorageGrafanaUsername,
			config.ValueTypeString,
			nil,
			"Username for basic authorization. No need if API key is set",
			true,
			"Grafana storage",
			nil,
			nil),
		config.NewVariable(
			annotations.ConfigStorageGrafanaPassword,
			config.ValueTypeString,
			nil,
			"Password for basic authorization. No need if API key is set",
			true,
			"Grafana storage",
			[]string{config.ViewPassword},
			nil),
		config.NewVariable(
			annotations.ConfigStorageGrafanaDashboards,
			config.ValueTypeString,
			nil,
			"Dashboards ID",
			true,
			"Grafana storage",
			[]string{config.ViewTags},
			map[string]interface{}{
				config.ViewOptionTagsDefaultText: "add a ID",
			}),
	}
}

func (c *Component) ConfigWatchers() []config.Watcher {
	return []config.Watcher{
		config.NewWatcher(annotations.ComponentName, []string{
			annotations.ConfigStorageGrafanaEnabled,
			annotations.ConfigStorageGrafanaAddress,
			annotations.ConfigStorageGrafanaApiKey,
			annotations.ConfigStorageGrafanaUsername,
			annotations.ConfigStorageGrafanaPassword,
		}, c.watchStorageGrafana),
	}
}

func (c *Component) watchStorageGrafana(_ string, _ interface{}, _ interface{}) {
	c.initStorageGrafana()
}
