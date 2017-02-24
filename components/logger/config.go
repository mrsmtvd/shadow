package logger

import (
	"github.com/kihamo/shadow/components/config"
)

const (
	ConfigLoggerLevel  = "logger.level"
	ConfigLoggerFields = "logger.fields"
)

func (c *Component) GetConfigVariables() []config.Variable {
	return []config.Variable{
		{
			Key:      ConfigLoggerLevel,
			Default:  6,
			Usage:    "Log level in RFC5424: 0 - Emergency, 1 - Alert, 2 - Critical, 3 - Error, 4 - Warning, 5 - Notice, 6 - Informational, 7 - Debug",
			Type:     config.ValueTypeInt,
			Editable: true,
		},
		{
			Key:      ConfigLoggerFields,
			Usage:    "Fields list with format: field_name=field1_value,field2_name=field2_value",
			Type:     config.ValueTypeString,
			Editable: true,
		},
	}
}

func (c *Component) GetConfigWatchers() map[string][]config.Watcher {
	return map[string][]config.Watcher{
		ConfigLoggerLevel:  {c.watchLoggerConfig},
		ConfigLoggerFields: {c.watchLoggerConfig},
	}
}

func (c *Component) watchLoggerConfig(_ string, newValue interface{}, _ interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.loggerConfig.Level = c.getLevel()
	c.loggerConfig.Fields = c.getDefaultFields()

	for key, _ := range c.loggers {
		c.loggers[key].(*logger).setConfig(c.loggerConfig)
	}
}
