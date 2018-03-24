package internal

import (
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/mail"
)

func (c *Component) ConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariable(
			mail.ConfigSmtpUsername,
			config.ValueTypeString,
			nil,
			"Username",
			true,
			"SMTP",
			nil,
			nil),
		config.NewVariable(
			mail.ConfigSmtpPassword,
			config.ValueTypeString,
			nil,
			"Password",
			true,
			"SMTP",
			[]string{config.ViewPassword},
			nil),
		config.NewVariable(
			mail.ConfigSmtpHost,
			config.ValueTypeString,
			nil,
			"Host",
			true,
			"SMTP",
			nil,
			nil),
		config.NewVariable(
			mail.ConfigSmtpPort,
			config.ValueTypeInt,
			25,
			"Port",
			true,
			"SMTP",
			nil,
			nil),
		config.NewVariable(
			mail.ConfigFromAddress,
			config.ValueTypeString,
			nil,
			"Mail from address",
			true,
			"Letter",
			nil,
			nil),
		config.NewVariable(
			mail.ConfigFromName,
			config.ValueTypeString,
			nil,
			"Mail from name",
			true,
			"Letter",
			nil,
			nil),
	}
}

func (c *Component) ConfigWatchers() []config.Watcher {
	return []config.Watcher{
		config.NewWatcher(mail.ComponentName, []string{
			mail.ConfigSmtpHost,
			mail.ConfigSmtpPort,
			mail.ConfigSmtpUsername,
			mail.ConfigSmtpPassword,
		}, c.watchForDialer),
	}
}

func (c *Component) watchForDialer(_ string, _ interface{}, _ interface{}) {
	c.initDialer(
		c.config.String(mail.ConfigSmtpHost),
		c.config.Int(mail.ConfigSmtpPort),
		c.config.String(mail.ConfigSmtpUsername),
		c.config.String(mail.ConfigSmtpPassword),
	)
}
