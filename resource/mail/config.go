package mail

import (
	"github.com/kihamo/shadow/resource/config"
)

func (r *Resource) GetConfigVariables() []config.Variable {
	return []config.Variable{
		{
			Key:   "mail.smtp.username",
			Value: "",
			Usage: "SMTP username",
		},
		{
			Key:   "mail.smtp.password",
			Value: "",
			Usage: "SMTP password",
		},
		{
			Key:   "mail.smtp.host",
			Value: "",
			Usage: "SMTP host",
		},
		{
			Key:   "mail.smtp.port",
			Value: 25,
			Usage: "SMTP port",
		},
		{
			Key:   "mail.from.address",
			Value: "",
			Usage: "Mail from address",
		},
		{
			Key:   "mail.from.name",
			Value: "",
			Usage: "Mail from name",
		},
	}
}
