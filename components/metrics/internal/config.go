package internal

import (
	"time"

	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/metrics"
)

func (c *Component) ConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariable(metrics.ConfigUrl, config.ValueTypeString).
			WithUsage("URL").
			WithGroup("InfluxDB").
			WithEditable(true),
		config.NewVariable(metrics.ConfigDatabase, config.ValueTypeString).
			WithUsage("Database name").
			WithGroup("InfluxDB").
			WithEditable(true),
		config.NewVariable(metrics.ConfigUsername, config.ValueTypeString).
			WithUsage("Username").
			WithGroup("InfluxDB").
			WithEditable(true),
		config.NewVariable(metrics.ConfigPassword, config.ValueTypeString).
			WithUsage("Password").
			WithGroup("InfluxDB").
			WithEditable(true).
			WithView([]string{config.ViewPassword}),
		config.NewVariable(metrics.ConfigPrecision, config.ValueTypeString).
			WithUsage("Precision").
			WithGroup("InfluxDB").
			WithEditable(true).
			WithDefault("s").
			WithView([]string{config.ViewEnum}).
			WithViewOptions(map[string]interface{}{
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
		config.NewVariable(metrics.ConfigInterval, config.ValueTypeDuration).
			WithUsage("Flush interval").
			WithEditable(true).
			WithDefault("1m"),
		config.NewVariable(metrics.ConfigLabels, config.ValueTypeString).
			WithUsage("Labels in format label1_name=label1_value").
			WithEditable(true).
			WithView([]string{config.ViewTags}).
			WithViewOptions(map[string]interface{}{config.ViewOptionTagsDefaultText: "add a label"}),
	}
}

func (c *Component) ConfigWatchers() []config.Watcher {
	return []config.Watcher{
		config.NewWatcher([]string{
			metrics.ConfigUrl,
			metrics.ConfigDatabase,
			metrics.ConfigUsername,
			metrics.ConfigPassword,
			metrics.ConfigPrecision,
		}, c.watchForStorage),
		config.NewWatcher([]string{metrics.ConfigInterval}, c.watchInterval),
		config.NewWatcher([]string{metrics.ConfigLabels}, c.watchLabels),
	}
}

func (c *Component) watchForStorage(_ string, newValue interface{}, _ interface{}) {
	_ = c.initStorage()
}

func (c *Component) watchInterval(_ string, newValue interface{}, _ interface{}) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	c.registry.SendInterval(newValue.(time.Duration))
}

func (c *Component) watchLabels(_ string, newValue interface{}, _ interface{}) {
	c.initLabels(newValue.(string))
}
