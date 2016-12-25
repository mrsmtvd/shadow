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
		ConfigMailSmtpHost:     {r.watchHost},
		ConfigMailSmtpPort:     {r.watchPort},
		ConfigMailSmtpUsername: {r.watchUsername},
		ConfigMailSmtpPassword: {r.watchPassword},
	}
}

func (r *Resource) watchHost(newValue interface{}, _ interface{}) {
	r.initDialer(
		newValue.(string),
		r.config.GetInt(ConfigMailSmtpPort),
		r.config.GetString(ConfigMailSmtpUsername),
		r.config.GetString(ConfigMailSmtpPassword),
	)
}

func (r *Resource) watchPort(newValue interface{}, _ interface{}) {
	r.initDialer(
		r.config.GetString(ConfigMailSmtpHost),
		newValue.(int),
		r.config.GetString(ConfigMailSmtpUsername),
		r.config.GetString(ConfigMailSmtpPassword),
	)
}

func (r *Resource) watchUsername(newValue interface{}, _ interface{}) {
	r.initDialer(
		r.config.GetString(ConfigMailSmtpHost),
		r.config.GetInt(ConfigMailSmtpPort),
		newValue.(string),
		r.config.GetString(ConfigMailSmtpPassword),
	)
}

func (r *Resource) watchPassword(newValue interface{}, _ interface{}) {
	r.initDialer(
		r.config.GetString(ConfigMailSmtpHost),
		r.config.GetInt(ConfigMailSmtpPort),
		r.config.GetString(ConfigMailSmtpUsername),
		newValue.(string),
	)
}
