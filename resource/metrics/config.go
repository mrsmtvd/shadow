package metrics

import "github.com/kihamo/shadow/resource/config"

func (r *Resource) GetConfigVariables() []config.Variable {
	return []config.Variable{
		{
			Key:   "metrics.url",
			Value: "",
			Usage: "InfluxDB url",
		},
		{
			Key:   "metrics.database",
			Value: "metrics",
			Usage: "InfluxDB database name",
		},
		{
			Key:   "metrics.username",
			Value: "",
			Usage: "InfluxDB username",
		},
		{
			Key:   "metrics.password",
			Value: "",
			Usage: "InfluxDB password",
		},
		{
			Key:   "metrics.interval",
			Value: "20s",
			Usage: "Flush interval",
		},
		{
			Key:   "metrics.tags",
			Value: "",
			Usage: "Tags list with format: tag1_name=tag1_value,tag2_name=tag2_value",
		},
		{
			Key:   "metrics.prefix",
			Value: "",
			Usage: "Prefix for measurements",
		},
	}
}
