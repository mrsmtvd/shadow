package internal

import (
	"github.com/mrsmtvd/shadow/components/config"
	"github.com/mrsmtvd/shadow/components/mail"
)

func (c *Component) ConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariable(mail.ConfigSMTPUsername, config.ValueTypeString).
			WithUsage("Username").
			WithGroup("SMTP").
			WithEditable(true),
		config.NewVariable(mail.ConfigSMTPPassword, config.ValueTypeString).
			WithUsage("Password").
			WithGroup("SMTP").
			WithEditable(true).
			WithView([]string{config.ViewPassword}),
		config.NewVariable(mail.ConfigSMTPHost, config.ValueTypeString).
			WithUsage("Host").
			WithGroup("SMTP").
			WithEditable(true),
		config.NewVariable(mail.ConfigSMTPPort, config.ValueTypeInt).
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
			mail.ConfigSMTPHost,
			mail.ConfigSMTPPort,
			mail.ConfigSMTPUsername,
			mail.ConfigSMTPPassword,
		}, c.watchForDialer),
	}
}

func (c *Component) watchForDialer(_ string, _ interface{}, _ interface{}) {
	c.initDialer(
		c.config.String(mail.ConfigSMTPHost),
		c.config.Int(mail.ConfigSMTPPort),
		c.config.String(mail.ConfigSMTPUsername),
		c.config.String(mail.ConfigSMTPPassword),
	)
}
