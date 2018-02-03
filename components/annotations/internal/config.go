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
			"Enabled Grafana storage",
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
			"Grafana ApiToken. No need if username and password is set",
			true,
			"Grafana storage",
			nil,
			nil),
		config.NewVariable(
			annotations.ConfigStorageGrafanaUsername,
			config.ValueTypeString,
			nil,
			"Grafana username for basic authorization. No need if api key is set",
			true,
			"Grafana storage",
			nil,
			nil),
		config.NewVariable(
			annotations.ConfigStorageGrafanaPassword,
			config.ValueTypeString,
			nil,
			"Grafana password for basic authorization. No need if api key is set",
			true,
			"Grafana storage",
			[]string{config.ViewPassword},
			nil),
		config.NewVariable(
			annotations.ConfigStorageGrafanaDashboards,
			config.ValueTypeString,
			nil,
			"Grafana dashboards id",
			true,
			"Grafana storage",
			[]string{config.ViewTags},
			map[string]interface{}{
				config.ViewOptionTagsDefaultText: "add a id",
			}),
		config.NewVariable(
			annotations.ConfigStorageTelegramEnabled,
			config.ValueTypeBool,
			false,
			"Enabled Telegram storage",
			true,
			"Telegram storage",
			nil,
			nil),
		config.NewVariable(
			annotations.ConfigStorageTelegramToken,
			config.ValueTypeString,
			nil,
			"Telegram bot token",
			true,
			"Telegram storage",
			nil,
			nil),
		config.NewVariable(
			annotations.ConfigStorageTelegramChats,
			config.ValueTypeString,
			nil,
			"Telegram chats id",
			true,
			"Telegram storage",
			[]string{config.ViewTags},
			map[string]interface{}{
				config.ViewOptionTagsDefaultText: "add a id",
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
		config.NewWatcher(annotations.ComponentName, []string{
			annotations.ConfigStorageTelegramEnabled,
			annotations.ConfigStorageTelegramToken,
			annotations.ConfigStorageTelegramChats,
		}, c.watchStorageTelegram),
	}
}

func (c *Component) watchStorageGrafana(_ string, _ interface{}, _ interface{}) {
	c.initStorageGrafana()
}

func (c *Component) watchStorageTelegram(_ string, _ interface{}, _ interface{}) {
	c.initStorageTelegram()
}
