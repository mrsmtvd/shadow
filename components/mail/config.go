package mail

import (
	"github.com/kihamo/shadow/components/config"
)

const (
	ConfigSmtpUsername = ComponentName + ".smtp.username"
	ConfigSmtpPassword = ComponentName + ".smtp.password"
	ConfigSmtpHost     = ComponentName + ".smtp.host"
	ConfigSmtpPort     = ComponentName + ".smtp.port"
	ConfigFromAddress  = ComponentName + ".from.address"
	ConfigFromName     = ComponentName + ".from.name"
)

func (c *Component) GetConfigVariables() []config.Variable {
	return []config.Variable{
		{
			Key:      ConfigSmtpUsername,
			Usage:    "SMTP username",
			Type:     config.ValueTypeString,
			Editable: true,
		},
		{
			Key:      ConfigSmtpPassword,
			Usage:    "SMTP password",
			Type:     config.ValueTypeString,
			Editable: true,
			View:     []string{config.ViewPassword},
		},
		{
			Key:      ConfigSmtpHost,
			Usage:    "SMTP host",
			Type:     config.ValueTypeString,
			Editable: true,
		},
		{
			Key:      ConfigSmtpPort,
			Default:  25,
			Usage:    "SMTP port",
			Type:     config.ValueTypeInt,
			Editable: true,
		},
		{
			Key:      ConfigFromAddress,
			Usage:    "Mail from address",
			Type:     config.ValueTypeString,
			Editable: true,
		},
		{
			Key:      ConfigFromName,
			Usage:    "Mail from name",
			Type:     config.ValueTypeString,
			Editable: true,
		},
	}
}

func (c *Component) GetConfigWatchers() map[string][]config.Watcher {
	return map[string][]config.Watcher{
		ConfigSmtpHost:     {c.watchHost},
		ConfigSmtpPort:     {c.watchPort},
		ConfigSmtpUsername: {c.watchUsername},
		ConfigSmtpPassword: {c.watchPassword},
	}
}

func (c *Component) watchHost(_ string, newValue interface{}, _ interface{}) {
	c.initDialer(
		newValue.(string),
		c.config.GetInt(ConfigSmtpPort),
		c.config.GetString(ConfigSmtpUsername),
		c.config.GetString(ConfigSmtpPassword),
	)
}

func (c *Component) watchPort(_ string, newValue interface{}, _ interface{}) {
	c.initDialer(
		c.config.GetString(ConfigSmtpHost),
		newValue.(int),
		c.config.GetString(ConfigSmtpUsername),
		c.config.GetString(ConfigSmtpPassword),
	)
}

func (c *Component) watchUsername(_ string, newValue interface{}, _ interface{}) {
	c.initDialer(
		c.config.GetString(ConfigSmtpHost),
		c.config.GetInt(ConfigSmtpPort),
		newValue.(string),
		c.config.GetString(ConfigSmtpPassword),
	)
}

func (c *Component) watchPassword(_ string, newValue interface{}, _ interface{}) {
	c.initDialer(
		c.config.GetString(ConfigSmtpHost),
		c.config.GetInt(ConfigSmtpPort),
		c.config.GetString(ConfigSmtpUsername),
		newValue.(string),
	)
}
