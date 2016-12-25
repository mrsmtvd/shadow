package logger

import (
	"github.com/kihamo/shadow/resource/config"
)

const (
	ConfigLoggerLevel  = "logger.level"
	ConfigLoggerFields = "logger.fields"
)

func (r *Resource) GetConfigVariables() []config.Variable {
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

func (r *Resource) GetConfigWatchers() map[string][]config.Watcher {
	return map[string][]config.Watcher{
		ConfigLoggerLevel:  {r.watchLoggerConfig},
		ConfigLoggerFields: {r.watchLoggerConfig},
	}
}

func (r *Resource) watchLoggerConfig(newValue interface{}, _ interface{}) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.loggerConfig.Level = r.getLevel()
	r.loggerConfig.Fields = r.getDefaultFields()

	for key, _ := range r.loggers {
		r.loggers[key].(*logger).setConfig(r.loggerConfig)
	}
}
