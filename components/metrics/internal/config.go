package internal

import (
	"time"

	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/metrics"
)

func (c *Component) GetConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariableItem(
			metrics.ConfigUrl,
			config.ValueTypeString,
			nil,
			"InfluxDB url",
			true,
			nil,
			nil),
		config.NewVariableItem(
			metrics.ConfigDatabase,
			config.ValueTypeString,
			nil,
			"InfluxDB database name",
			true,
			nil,
			nil),
		config.NewVariableItem(
			metrics.ConfigUsername,
			config.ValueTypeString,
			nil,
			"InfluxDB username",
			true,
			nil,
			nil),
		config.NewVariableItem(
			metrics.ConfigPassword,
			config.ValueTypeString,
			nil,
			"InfluxDB password",
			true,
			[]string{config.ViewPassword},
			nil),
		config.NewVariableItem(
			metrics.ConfigPrecision,
			config.ValueTypeString,
			"s",
			"InfluxDB precision",
			true,
			nil,
			nil),
		config.NewVariableItem(
			metrics.ConfigInterval,
			config.ValueTypeDuration,
			"1m",
			"Flush interval",
			true,
			nil,
			nil),
		config.NewVariableItem(
			metrics.ConfigLabels,
			config.ValueTypeString,
			nil,
			"Labels list with format: label1_name=label1_value",
			true,
			[]string{config.ViewTags},
			map[string]interface{}{
				config.ViewOptionTagsDefaultText: "add a label",
			}),
	}
}

func (c *Component) GetConfigWatchers() []config.Watcher {
	return []config.Watcher{
		config.NewWatcher(metrics.ComponentName, []string{
			metrics.ConfigUrl,
			metrics.ConfigDatabase,
			metrics.ConfigUsername,
			metrics.ConfigPassword,
			metrics.ConfigPrecision,
		}, c.watchForStorage),
		config.NewWatcher(metrics.ComponentName, []string{metrics.ConfigInterval}, c.watchInterval),
		config.NewWatcher(metrics.ComponentName, []string{metrics.ConfigLabels}, c.watchLabels),
	}
}

func (c *Component) watchForStorage(_ string, newValue interface{}, _ interface{}) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	c.storage.Reinitialization(
		c.config.GetString(metrics.ConfigUrl),
		c.config.GetString(metrics.ConfigDatabase),
		c.config.GetString(metrics.ConfigUsername),
		c.config.GetString(metrics.ConfigPassword),
		c.config.GetString(metrics.ConfigPrecision))
}

func (c *Component) watchInterval(_ string, newValue interface{}, _ interface{}) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	c.registry.SendInterval(newValue.(time.Duration))
}

func (c *Component) watchLabels(_ string, newValue interface{}, _ interface{}) {
	c.initLabels(newValue.(string))
}
