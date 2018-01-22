package internal

import (
	"github.com/kihamo/shadow/components/annotations"
	"github.com/kihamo/shadow/components/config"
)

func (c *Component) GetConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariable(
			annotations.ConfigStorageGrafanaAddress,
			config.ValueTypeString,
			nil,
			"Grafana address of HTTP API in format http://host:port",
			true,
			nil,
			nil),
		config.NewVariable(
			annotations.ConfigStorageGrafanaApiKey,
			config.ValueTypeString,
			nil,
			"Grafana ApiToken. No need if username and password is set",
			true,
			nil,
			nil),
		config.NewVariable(
			annotations.ConfigStorageGrafanaUsername,
			config.ValueTypeString,
			nil,
			"Grafana username for basic authorization. No need if api key is set",
			true,
			nil,
			nil),
		config.NewVariable(
			annotations.ConfigStorageGrafanaPassword,
			config.ValueTypeString,
			nil,
			"Grafana password for basic authorization. No need if api key is set",
			true,
			[]string{config.ViewPassword},
			nil),
		config.NewVariable(
			annotations.ConfigStorageGrafanaDashboards,
			config.ValueTypeString,
			nil,
			"Grafana dashboards id",
			true,
			[]string{config.ViewTags},
			map[string]interface{}{
				config.ViewOptionTagsDefaultText: "add a id",
			}),
	}
}

func (c *Component) GetConfigWatchers() []config.Watcher {
	return []config.Watcher{
		config.NewWatcher(annotations.ComponentName, []string{
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
