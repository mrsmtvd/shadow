package metrics

import (
	"github.com/kihamo/shadow/resource/config"
)

func (r *Resource) GetConfigVariables() []config.Variable {
	return []config.Variable{
		{
			Key:   "metrics.url",
			Usage: "InfluxDB url",
			Type:  config.ValueTypeString,
		},
		{
			Key:   "metrics.database",
			Usage: "InfluxDB database name",
			Type:  config.ValueTypeString,
		},
		{
			Key:   "metrics.username",
			Usage: "InfluxDB username",
			Type:  config.ValueTypeString,
		},
		{
			Key:   "metrics.password",
			Usage: "InfluxDB password",
			Type:  config.ValueTypeString,
		},
		{
			Key:     "metrics.interval",
			Default: "30s",
			Usage:   "Flush interval",
			Type:    config.ValueTypeDuration,
		},
		{
			Key:   "metrics.tags",
			Usage: "Tags list with format: tag1_name=tag1_value,tag2_name=tag2_value",
			Type:  config.ValueTypeString,
		},
		{
			Key:   "metrics.prefix",
			Usage: "Prefix for measurements",
			Type:  config.ValueTypeString,
		},
	}
}
