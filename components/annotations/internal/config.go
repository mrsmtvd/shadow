package internal

import (
	"github.com/kihamo/shadow/components/annotations"
	"github.com/kihamo/shadow/components/config"
)

func (c *Component) ConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariable(annotations.ConfigStorageGrafanaEnabled, config.ValueTypeBool).
			WithUsage("Enabled").
			WithGroup("Grafana storage").
			WithEditable(true),
		config.NewVariable(annotations.ConfigStorageGrafanaAddress, config.ValueTypeString).
			WithUsage("Grafana address of HTTP API in format http://host:port").
			WithGroup("Grafana storage").
			WithEditable(true),
		config.NewVariable(annotations.ConfigStorageGrafanaAPIKey, config.ValueTypeString).
			WithUsage("API key. No need if username and password is set").
			WithGroup("Grafana storage").
			WithEditable(true),
		config.NewVariable(annotations.ConfigStorageGrafanaUsername, config.ValueTypeString).
			WithUsage("Username for basic authorization. No need if API key is set").
			WithGroup("Grafana storage").
			WithEditable(true),
		config.NewVariable(annotations.ConfigStorageGrafanaPassword, config.ValueTypeString).
			WithUsage("Password for basic authorization. No need if API key is set").
			WithGroup("Grafana storage").
			WithEditable(true).
			WithView([]string{config.ViewPassword}),
		config.NewVariable(annotations.ConfigStorageGrafanaDashboards, config.ValueTypeString).
			WithUsage("Dashboards ID").
			WithGroup("Grafana storage").
			WithEditable(true).
			WithView([]string{config.ViewTags}).
			WithViewOptions(map[string]interface{}{config.ViewOptionTagsDefaultText: "add a ID"}),
	}
}

func (c *Component) ConfigWatchers() []config.Watcher {
	return []config.Watcher{
		config.NewWatcher([]string{
			annotations.ConfigStorageGrafanaEnabled,
			annotations.ConfigStorageGrafanaAddress,
			annotations.ConfigStorageGrafanaAPIKey,
			annotations.ConfigStorageGrafanaUsername,
			annotations.ConfigStorageGrafanaPassword,
		}, c.watchStorageGrafana),
	}
}

func (c *Component) watchStorageGrafana(_ string, _ interface{}, _ interface{}) {
	c.initStorageGrafana()
}
