package internal

import (
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/logging"
	"github.com/kihamo/shadow/components/logging/output"
)

func (c *Component) ConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariable(logging.ConfigLevel, config.ValueTypeInt).
			WithUsage("Log level in format RFC5424").
			WithEditable(true).
			WithDefault(logging.LevelInformational).
			WithView([]string{config.ViewEnum}).
			WithViewOptions(map[string]interface{}{
				config.ViewOptionEnumOptions: [][]interface{}{
					{logging.LevelEmergency, "Emergency"},
					{logging.LevelAlert, "Alert"},
					{logging.LevelCritical, "Critical"},
					{logging.LevelError, "Error"},
					{logging.LevelWarning, "Warning"},
					{logging.LevelNotice, "Notice"},
					{logging.LevelInformational, "Informational"},
					{logging.LevelDebug, "Debug"},
				},
			}),
		config.NewVariable(logging.ConfigFields, config.ValueTypeString).
			WithUsage("Fields in format field_name=field1_value,field2_name=field2_value").
			WithEditable(true).
			WithView([]string{config.ViewTags}).
			WithViewOptions(map[string]interface{}{config.ViewOptionTagsDefaultText: "add a field"}),
	}
}

func (c *Component) ConfigWatchers() []config.Watcher {
	return []config.Watcher{
		config.NewWatcher([]string{logging.ConfigLevel}, c.watchLoggerLevel),
		config.NewWatcher([]string{logging.ConfigFields}, c.watchLoggerFields),
	}
}

func (c *Component) watchLoggerLevel(_ string, newValue interface{}, _ interface{}) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	level := output.ConvertLoggerToXLogLevel(c.getLevel())

	for key := range c.loggers {
		c.loggers[key].(*output.WrapperXLog).SetLevel(level)
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
		l := c.loggers[key].(*output.WrapperXLog)

		existsFields := map[string]interface{}{}

		for k, v := range l.GetFields() {
			if _, ok := removeFields[k]; !ok {
				existsFields[k] = v
			}
		}

		for k, v := range globalFields {
			existsFields[k] = v
		}

		l.SetFields(existsFields)
	}
}
