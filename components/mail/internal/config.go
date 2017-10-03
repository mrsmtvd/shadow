package internal

import (
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/mail"
)

func (c *Component) GetConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariableItem(
			mail.ConfigSmtpUsername,
			config.ValueTypeString,
			nil,
			"SMTP username",
			true,
			nil,
			nil),
		config.NewVariableItem(
			mail.ConfigSmtpPassword,
			config.ValueTypeString,
			nil,
			"SMTP password",
			true,
			[]string{config.ViewPassword},
			nil),
		config.NewVariableItem(
			mail.ConfigSmtpHost,
			config.ValueTypeString,
			nil,
			"SMTP host",
			true,
			nil,
			nil),
		config.NewVariableItem(
			mail.ConfigSmtpPort,
			config.ValueTypeInt,
			25,
			"SMTP port",
			true,
			nil,
			nil),
		config.NewVariableItem(
			mail.ConfigFromAddress,
			config.ValueTypeString,
			nil,
			"Mail from address",
			true,
			nil,
			nil),
		config.NewVariableItem(
			mail.ConfigFromName,
			config.ValueTypeString,
			nil,
			"Mail from name",
			true,
			nil,
			nil),
	}
}

func (c *Component) GetConfigWatchers() map[string][]config.Watcher {
	return map[string][]config.Watcher{
		mail.ConfigSmtpHost:     {c.watchHost},
		mail.ConfigSmtpPort:     {c.watchPort},
		mail.ConfigSmtpUsername: {c.watchUsername},
		mail.ConfigSmtpPassword: {c.watchPassword},
	}
}

func (c *Component) watchHost(_ string, newValue interface{}, _ interface{}) {
	c.initDialer(
		newValue.(string),
		c.config.GetInt(mail.ConfigSmtpPort),
		c.config.GetString(mail.ConfigSmtpUsername),
		c.config.GetString(mail.ConfigSmtpPassword),
	)
}

func (c *Component) watchPort(_ string, newValue interface{}, _ interface{}) {
	c.initDialer(
		c.config.GetString(mail.ConfigSmtpHost),
		newValue.(int),
		c.config.GetString(mail.ConfigSmtpUsername),
		c.config.GetString(mail.ConfigSmtpPassword),
	)
}

func (c *Component) watchUsername(_ string, newValue interface{}, _ interface{}) {
	c.initDialer(
		c.config.GetString(mail.ConfigSmtpHost),
		c.config.GetInt(mail.ConfigSmtpPort),
		newValue.(string),
		c.config.GetString(mail.ConfigSmtpPassword),
	)
}

func (c *Component) watchPassword(_ string, newValue interface{}, _ interface{}) {
	c.initDialer(
		c.config.GetString(mail.ConfigSmtpHost),
		c.config.GetInt(mail.ConfigSmtpPort),
		c.config.GetString(mail.ConfigSmtpUsername),
		newValue.(string),
	)
}
