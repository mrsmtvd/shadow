package mail

import (
	"github.com/kihamo/shadow/resource/config"
)

const (
	ConfigMailSmtpUsername = "mail.smtp.username"
	ConfigMailSmtpPassword = "mail.smtp.password"
	ConfigMailSmtpHost     = "mail.smtp.host"
	ConfigMailSmtpPort     = "mail.smtp.port"
	ConfigMailFromAddress  = "mail.from.address"
	ConfigMailFromName     = "mail.from.name"
)

func (r *Resource) GetConfigVariables() []config.Variable {
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

func (r *Resource) GetConfigWatchers() map[string][]config.Watcher {
	return map[string][]config.Watcher{
		ConfigMailSmtpUsername: {r.watchSmtp},
		ConfigMailSmtpPassword: {r.watchSmtp},
		ConfigMailSmtpHost:     {r.watchSmtp},
		ConfigMailSmtpPort:     {r.watchSmtp},
	}
}

func (r *Resource) watchSmtp(_ interface{}, _ interface{}) {
	r.initDialer()
}
