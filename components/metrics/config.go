package metrics

import (
	"time"

	"github.com/kihamo/shadow/components/config"
)

const (
	ConfigUrl       = ComponentName + ".url"
	ConfigDatabase  = ComponentName + ".database"
	ConfigUsername  = ComponentName + ".username"
	ConfigPassword  = ComponentName + ".password"
	ConfigPrecision = ComponentName + ".precision"
	ConfigInterval  = ComponentName + ".interval"
	ConfigLabels    = ComponentName + ".labels"
)

func (c *Component) GetConfigVariables() []config.Variable {
	return []config.Variable{
		{
			Key:      ConfigUrl,
			Usage:    "InfluxDB url",
			Type:     config.ValueTypeString,
			Editable: true,
		},
		{
			Key:      ConfigDatabase,
			Usage:    "InfluxDB database name",
			Type:     config.ValueTypeString,
			Editable: true,
		},
		{
			Key:      ConfigUsername,
			Usage:    "InfluxDB username",
			Type:     config.ValueTypeString,
			Editable: true,
		},
		{
			Key:      ConfigPassword,
			Usage:    "InfluxDB password",
			Type:     config.ValueTypeString,
			Editable: true,
		},
		{
			Key:      ConfigPrecision,
			Usage:    "InfluxDB precision",
			Default:  "s",
			Type:     config.ValueTypeString,
			Editable: true,
		},
		{
			Key:      ConfigInterval,
			Default:  "1m",
			Usage:    "Flush interval",
			Type:     config.ValueTypeDuration,
			Editable: true,
		},
		{
			Key:      ConfigLabels,
			Usage:    "Labels list with format: label1_name=label1_value,label2_name=label2_value",
			Type:     config.ValueTypeString,
			Editable: true,
		},
	}
}

func (c *Component) GetConfigWatchers() map[string][]config.Watcher {
	return map[string][]config.Watcher{
		ConfigUrl:       {c.watchUrl},
		ConfigDatabase:  {c.watchDatabase},
		ConfigUsername:  {c.watchUsername},
		ConfigPassword:  {c.watchPassword},
		ConfigPrecision: {c.watchPrecision},
		ConfigInterval:  {c.watchInterval},
		ConfigLabels:    {c.watchLabels},
	}
}

func (c *Component) watchUrl(_ string, newValue interface{}, _ interface{}) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	c.storage.Reinitialization(
		newValue.(string),
		c.config.GetString(ConfigDatabase),
		c.config.GetString(ConfigUsername),
		c.config.GetString(ConfigPassword),
		c.config.GetString(ConfigPrecision))
}

func (c *Component) watchDatabase(_ string, newValue interface{}, _ interface{}) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	c.storage.Reinitialization(
		c.config.GetString(ConfigUrl),
		newValue.(string),
		c.config.GetString(ConfigUsername),
		c.config.GetString(ConfigPassword),
		c.config.GetString(ConfigPrecision))
}

func (c *Component) watchUsername(_ string, newValue interface{}, _ interface{}) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	c.storage.Reinitialization(
		c.config.GetString(ConfigUrl),
		c.config.GetString(ConfigDatabase),
		newValue.(string),
		c.config.GetString(ConfigPassword),
		c.config.GetString(ConfigPrecision))
}

func (c *Component) watchPassword(_ string, newValue interface{}, _ interface{}) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	c.storage.Reinitialization(
		c.config.GetString(ConfigUrl),
		c.config.GetString(ConfigDatabase),
		c.config.GetString(ConfigUsername),
		newValue.(string),
		c.config.GetString(ConfigPrecision))
}

func (c *Component) watchPrecision(_ string, newValue interface{}, _ interface{}) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	c.storage.Reinitialization(
		c.config.GetString(ConfigUrl),
		c.config.GetString(ConfigDatabase),
		c.config.GetString(ConfigUsername),
		c.config.GetString(ConfigPassword),
		newValue.(string))
}

func (c *Component) watchInterval(_ string, newValue interface{}, _ interface{}) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	c.registry.SendInterval(newValue.(time.Duration))
}

func (c *Component) watchLabels(_ string, newValue interface{}, _ interface{}) {
	c.initLabels(newValue.(string))
}
