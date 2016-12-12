package metrics

import (
	"github.com/kihamo/shadow/resource/config"
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

func (r *Resource) GetConfigVariables() []config.Variable {
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
			Key:     ConfigMetricsInterval,
			Default: "30s",
			Usage:   "Flush interval",
			Type:    config.ValueTypeDuration,
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

func (r *Resource) GetConfigWatchers() map[string][]config.Watcher {
	return map[string][]config.Watcher{
		ConfigMetricsUrl:      {r.watchClient},
		ConfigMetricsUsername: {r.watchClient},
		ConfigMetricsPassword: {r.watchClient},
	}
}

func (r *Resource) watchClient(_ interface{}, _ interface{}) {
	r.initClient()
}
