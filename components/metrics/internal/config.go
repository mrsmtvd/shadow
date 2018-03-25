package internal

import (
	"time"

	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/metrics"
)

func (c *Component) ConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariable(
			metrics.ConfigUrl,
			config.ValueTypeString,
			nil,
			"URL",
			true,
			"InfluxDB",
			nil,
			nil),
		config.NewVariable(
			metrics.ConfigDatabase,
			config.ValueTypeString,
			nil,
			"Database name",
			true,
			"InfluxDB",
			nil,
			nil),
		config.NewVariable(
			metrics.ConfigUsername,
			config.ValueTypeString,
			nil,
			"Username",
			true,
			"InfluxDB",
			nil,
			nil),
		config.NewVariable(
			metrics.ConfigPassword,
			config.ValueTypeString,
			nil,
			"Password",
			true,
			"InfluxDB",
			[]string{config.ViewPassword},
			nil),
		config.NewVariable(
			metrics.ConfigPrecision,
			config.ValueTypeString,
			"s",
			"Precision",
			true,
			"InfluxDB",
			[]string{
				config.ViewEnum,
			},
			map[string]interface{}{
				config.ViewOptionEnumOptions: [][]interface{}{
					{"rfc3339", "rfc3339"},
					{"h", "Hour"},
					{"m", "Minute"},
					{"s", "Second"},
					{"ms", "Millisecond"},
					{"u", "Microsecond"},
					{"ns", "Nanosecond"},
				},
			}),
		config.NewVariable(
			metrics.ConfigInterval,
			config.ValueTypeDuration,
			"1m",
			"Flush interval",
			true,
			"Others",
			nil,
			nil),
		config.NewVariable(
			metrics.ConfigLabels,
			config.ValueTypeString,
			nil,
			"Labels in format label1_name=label1_value",
			true,
			"Others",
			[]string{config.ViewTags},
			map[string]interface{}{
				config.ViewOptionTagsDefaultText: "add a label",
			}),
	}
}

func (c *Component) ConfigWatchers() []config.Watcher {
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
		c.config.String(metrics.ConfigUrl),
		c.config.String(metrics.ConfigDatabase),
		c.config.String(metrics.ConfigUsername),
		c.config.String(metrics.ConfigPassword),
		c.config.String(metrics.ConfigPrecision))
}

func (c *Component) watchInterval(_ string, newValue interface{}, _ interface{}) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	c.registry.SendInterval(newValue.(time.Duration))
}

func (c *Component) watchLabels(_ string, newValue interface{}, _ interface{}) {
	c.initLabels(newValue.(string))
}
