package mail

import (
	"github.com/kihamo/shadow/resource/config"
)

func (r *Resource) GetConfigVariables() []config.Variable {
	return []config.Variable{
		{
			Key:   "mail.smtp.username",
			Usage: "SMTP username",
			Type:  config.ValueTypeString,
		},
		{
			Key:   "mail.smtp.password",
			Usage: "SMTP password",
			Type:  config.ValueTypeString,
		},
		{
			Key:   "mail.smtp.host",
			Usage: "SMTP host",
			Type:  config.ValueTypeString,
		},
		{
			Key:     "mail.smtp.port",
			Default: 25,
			Usage:   "SMTP port",
			Type:    config.ValueTypeInt,
		},
		{
			Key:   "mail.from.address",
			Usage: "Mail from address",
			Type:  config.ValueTypeString,
		},
		{
			Key:   "mail.from.name",
			Usage: "Mail from name",
			Type:  config.ValueTypeString,
		},
	}
}
