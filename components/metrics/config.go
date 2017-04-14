package metrics

import (
	"time"

	"github.com/kihamo/shadow/components/config"
)

const (
	ConfigMetricsUrl       = ComponentName + ".url"
	ConfigMetricsDatabase  = ComponentName + ".database"
	ConfigMetricsUsername  = ComponentName + ".username"
	ConfigMetricsPassword  = ComponentName + ".password"
	ConfigMetricsPrecision = ComponentName + ".precision"
	ConfigMetricsInterval  = ComponentName + ".interval"
	ConfigMetricsLabels    = ComponentName + ".labels"
	ConfigMetricsPrefix    = ComponentName + ".prefix"
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
			Key:      ConfigMetricsDatabase,
			Usage:    "InfluxDB database name",
			Type:     config.ValueTypeString,
			Editable: true,
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
			Key:      ConfigMetricsPrecision,
			Usage:    "InfluxDB precision",
			Default:  "s",
			Type:     config.ValueTypeString,
			Editable: true,
		},
		{
			Key:      ConfigMetricsInterval,
			Default:  "1m",
			Usage:    "Flush interval",
			Type:     config.ValueTypeDuration,
			Editable: true,
		},
		{
			Key:      ConfigMetricsLabels,
			Usage:    "Labels list with format: label1_name=label1_value,label2_name=label2_value",
			Type:     config.ValueTypeString,
			Editable: true,
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
		ConfigMetricsUrl:       {c.watchUrl},
		ConfigMetricsDatabase:  {c.watchDatabase},
		ConfigMetricsUsername:  {c.watchUsername},
		ConfigMetricsPassword:  {c.watchPassword},
		ConfigMetricsPrecision: {c.watchPrecision},
		ConfigMetricsInterval:  {c.watchInterval},
		ConfigMetricsLabels:    {c.watchLabels},
	}
}

func (c *Component) watchUrl(_ string, newValue interface{}, _ interface{}) {
	/*
		c.initCollector(
			newValue.(string),
			c.config.GetString(ConfigMetricsDatabase),
			c.config.GetString(ConfigMetricsUsername),
			c.config.GetString(ConfigMetricsPassword),
			c.config.GetString(ConfigMetricsPrecision),
		)
	*/
}

func (c *Component) watchDatabase(_ string, newValue interface{}, _ interface{}) {
	/*
		c.initCollector(
			c.config.GetString(ConfigMetricsUrl),
			newValue.(string),
			c.config.GetString(ConfigMetricsUsername),
			c.config.GetString(ConfigMetricsPassword),
			c.config.GetString(ConfigMetricsPrecision),
		)
	*/
}

func (c *Component) watchUsername(_ string, newValue interface{}, _ interface{}) {
	/*
		c.initCollector(
			c.config.GetString(ConfigMetricsUrl),
			c.config.GetString(ConfigMetricsDatabase),
			newValue.(string),
			c.config.GetString(ConfigMetricsPassword),
			c.config.GetString(ConfigMetricsPrecision),
		)
	*/
}

func (c *Component) watchPassword(_ string, newValue interface{}, _ interface{}) {
	/*
		c.initCollector(
			c.config.GetString(ConfigMetricsUrl),
			c.config.GetString(ConfigMetricsDatabase),
			c.config.GetString(ConfigMetricsUsername),
			newValue.(string),
			c.config.GetString(ConfigMetricsPrecision),
		)
	*/
}

func (c *Component) watchPrecision(_ string, newValue interface{}, _ interface{}) {
	/*
		c.initCollector(
			c.config.GetString(ConfigMetricsUrl),
			c.config.GetString(ConfigMetricsDatabase),
			c.config.GetString(ConfigMetricsUsername),
			c.config.GetString(ConfigMetricsPassword),
			newValue.(string),
		)
	*/
}

func (c *Component) watchInterval(_ string, newValue interface{}, _ interface{}) {
	c.changeTicker <- newValue.(time.Duration)
}

func (c *Component) watchLabels(_ string, newValue interface{}, _ interface{}) {
	//c.initLabels(newValue.(string))
}
