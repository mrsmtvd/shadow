package metrics

import (
	"time"

	"github.com/kihamo/shadow/components/config"
)

const (
	ConfigMetricsUrl       = "metrics.url"
	ConfigMetricsDatabase  = "metrics.database"
	ConfigMetricsUsername  = "metrics.username"
	ConfigMetricsPassword  = "metrics.password"
	ConfigMetricsPrecision = "metrics.precision"
	ConfigMetricsInterval  = "metrics.interval"
	ConfigMetricsTags      = "metrics.tags"
	ConfigMetricsPrefix    = "metrics.prefix"
)

func (c *Component) GetConfigVariables() []config.Variable {
	return []config.Variable{
		{
			Key:      ConfigMetricsUrl,
			Usage:    "InfluxDB url",
			Type:     config.ValueTypeString,
			Editable: true,
		},
		{
			Key:   ConfigMetricsDatabase,
			Usage: "InfluxDB database name",
			Type:  config.ValueTypeString,
		},
		{
			Key:      ConfigMetricsUsername,
			Usage:    "InfluxDB username",
			Type:     config.ValueTypeString,
			Editable: true,
		},
		{
			Key:      ConfigMetricsPassword,
			Usage:    "InfluxDB password",
			Type:     config.ValueTypeString,
			Editable: true,
		},
		{
			Key:     ConfigMetricsPrecision,
			Usage:   "InfluxDB precision",
			Default: "s",
			Type:    config.ValueTypeString,
		},
		{
			Key:      ConfigMetricsInterval,
			Default:  "1m",
			Usage:    "Flush interval",
			Type:     config.ValueTypeDuration,
			Editable: true,
		},
		{
			Key:   ConfigMetricsTags,
			Usage: "Tags list with format: tag1_name=tag1_value,tag2_name=tag2_value",
			Type:  config.ValueTypeString,
		},
		{
			Key:   ConfigMetricsPrefix,
			Usage: "Prefix for measurements",
			Type:  config.ValueTypeString,
		},
	}
}

func (c *Component) GetConfigWatchers() map[string][]config.Watcher {
	return map[string][]config.Watcher{
		ConfigMetricsUrl:      {c.watchUrl},
		ConfigMetricsUsername: {c.watchUsername},
		ConfigMetricsPassword: {c.watchPassword},
		ConfigMetricsInterval: {c.watchInterval},
	}
}

func (c *Component) watchUrl(_ string, newValue interface{}, _ interface{}) {
	c.initClient(
		newValue.(string),
		c.config.GetString(ConfigMetricsUsername),
		c.config.GetString(ConfigMetricsPassword),
	)
}

func (c *Component) watchUsername(_ string, newValue interface{}, _ interface{}) {
	c.initClient(
		c.config.GetString(ConfigMetricsUrl),
		newValue.(string),
		c.config.GetString(ConfigMetricsPassword),
	)
}

func (c *Component) watchPassword(_ string, newValue interface{}, _ interface{}) {
	c.initClient(
		c.config.GetString(ConfigMetricsUrl),
		c.config.GetString(ConfigMetricsUsername),
		newValue.(string),
	)
}

func (c *Component) watchInterval(_ string, newValue interface{}, _ interface{}) {
	c.changeTicker <- newValue.(time.Duration)
}
