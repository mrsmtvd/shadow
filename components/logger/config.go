package logger

import (
	"github.com/kihamo/shadow/components/config"
)

const (
	ConfigLevel  = ComponentName + ".level"
	ConfigFields = ComponentName + ".fields"
)

func (c *Component) GetConfigVariables() []config.Variable {
	return []config.Variable{
		{
			Key:      ConfigLevel,
			Default:  6,
			Usage:    "Log level in RFC5424: 0 - Emergency, 1 - Alert, 2 - Critical, 3 - Error, 4 - Warning, 5 - Notice, 6 - Informational, 7 - Debug",
			Type:     config.ValueTypeInt,
			Editable: true,
		},
		{
			Key:      ConfigFields,
			Usage:    "Fields list with format: field_name=field1_value,field2_name=field2_value",
			Type:     config.ValueTypeString,
			Editable: true,
			View:     []string{config.ViewTags},
			ViewOptions: map[string]interface{}{
				config.ViewOptionTagsDefaultText: "add a field",
			},
		},
	}
}

func (c *Component) GetConfigWatchers() map[string][]config.Watcher {
	return map[string][]config.Watcher{
		ConfigLevel:  {c.watchLoggerLevel},
		ConfigFields: {c.watchLoggerFields},
	}
}

func (c *Component) watchLoggerLevel(_ string, newValue interface{}, _ interface{}) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	level := c.getConfigLevel()

	for key, _ := range c.loggers {
		c.loggers[key].(*logger).setLevel(level)
	}
}

func (c *Component) watchLoggerFields(_ string, newValue interface{}, oldValue interface{}) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	removeFields := map[string]struct{}{}
	newFields := c.parseFields(newValue.(string))
	oldFields := c.parseFields(oldValue.(string))

	for k, _ := range oldFields {
		if _, ok := newFields[k]; !ok {
			removeFields[k] = struct{}{}
		}
	}

	globalFields := c.getFields()

	for key, _ := range c.loggers {
		l := c.loggers[key].(*logger)

		existsFields := map[string]interface{}{}

		for k, v := range l.GetFields() {
			if _, ok := removeFields[k]; !ok {
				existsFields[k] = v
			}
		}

		for k, v := range globalFields {
			existsFields[k] = v
		}

		l.setFields(existsFields)
	}
}
