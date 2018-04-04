package internal

import (
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/mail"
)

func (c *Component) ConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariable(mail.ConfigSmtpUsername, config.ValueTypeString).
			WithUsage("Username").
			WithGroup("SMTP").
			WithEditable(true),
		config.NewVariable(mail.ConfigSmtpPassword, config.ValueTypeString).
			WithUsage("Password").
			WithGroup("SMTP").
			WithEditable(true).
			WithView([]string{config.ViewPassword}),
		config.NewVariable(mail.ConfigSmtpHost, config.ValueTypeString).
			WithUsage("Host").
			WithGroup("SMTP").
			WithEditable(true),
		config.NewVariable(mail.ConfigSmtpPort, config.ValueTypeInt).
			WithUsage("Port").
			WithGroup("SMTP").
			WithEditable(true).
			WithDefault(25),
		config.NewVariable(mail.ConfigFromAddress, config.ValueTypeString).
			WithUsage("Mail from address").
			WithGroup("Letter").
			WithEditable(true),
		config.NewVariable(mail.ConfigFromName, config.ValueTypeString).
			WithUsage("Mail from name").
			WithGroup("Letter").
			WithEditable(true),
	}
}

func (c *Component) ConfigWatchers() []config.Watcher {
	return []config.Watcher{
		config.NewWatcher([]string{
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
