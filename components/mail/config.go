package mail

import (
	"github.com/kihamo/shadow/components/config"
)

const (
	ConfigMailSmtpUsername = "mail.smtp.username"
	ConfigMailSmtpPassword = "mail.smtp.password"
	ConfigMailSmtpHost     = "mail.smtp.host"
	ConfigMailSmtpPort     = "mail.smtp.port"
	ConfigMailFromAddress  = "mail.from.address"
	ConfigMailFromName     = "mail.from.name"
)

func (c *Component) GetConfigVariables() []config.Variable {
	return []config.Variable{
		{
			Key:      ConfigMailSmtpUsername,
			Usage:    "SMTP username",
			Type:     config.ValueTypeString,
			Editable: true,
		},
		{
			Key:      ConfigMailSmtpPassword,
			Usage:    "SMTP password",
			Type:     config.ValueTypeString,
			Editable: true,
		},
		{
			Key:      ConfigMailSmtpHost,
			Usage:    "SMTP host",
			Type:     config.ValueTypeString,
			Editable: true,
		},
		{
			Key:      ConfigMailSmtpPort,
			Default:  25,
			Usage:    "SMTP port",
			Type:     config.ValueTypeInt,
			Editable: true,
		},
		{
			Key:      ConfigMailFromAddress,
			Usage:    "Mail from address",
			Type:     config.ValueTypeString,
			Editable: true,
		},
		{
			Key:      ConfigMailFromName,
			Usage:    "Mail from name",
			Type:     config.ValueTypeString,
			Editable: true,
		},
	}
}

func (c *Component) GetConfigWatchers() map[string][]config.Watcher {
	return map[string][]config.Watcher{
		ConfigMailSmtpHost:     {c.watchHost},
		ConfigMailSmtpPort:     {c.watchPort},
		ConfigMailSmtpUsername: {c.watchUsername},
		ConfigMailSmtpPassword: {c.watchPassword},
	}
}

func (c *Component) watchHost(newValue interface{}, _ interface{}) {
	c.initDialer(
		newValue.(string),
		c.config.GetInt(ConfigMailSmtpPort),
		c.config.GetString(ConfigMailSmtpUsername),
		c.config.GetString(ConfigMailSmtpPassword),
	)
}

func (c *Component) watchPort(newValue interface{}, _ interface{}) {
	c.initDialer(
		c.config.GetString(ConfigMailSmtpHost),
		newValue.(int),
		c.config.GetString(ConfigMailSmtpUsername),
		c.config.GetString(ConfigMailSmtpPassword),
	)
}

func (c *Component) watchUsername(newValue interface{}, _ interface{}) {
	c.initDialer(
		c.config.GetString(ConfigMailSmtpHost),
		c.config.GetInt(ConfigMailSmtpPort),
		newValue.(string),
		c.config.GetString(ConfigMailSmtpPassword),
	)
}

func (c *Component) watchPassword(newValue interface{}, _ interface{}) {
	c.initDialer(
		c.config.GetString(ConfigMailSmtpHost),
		c.config.GetInt(ConfigMailSmtpPort),
		c.config.GetString(ConfigMailSmtpUsername),
		newValue.(string),
	)
}
