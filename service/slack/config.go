package slack

import (
	"github.com/kihamo/shadow/resource"
)

func (s *SlackService) GetConfigVariables() []resource.ConfigVariable {
	return []resource.ConfigVariable{
		resource.ConfigVariable{
			Key:   "slack-token",
			Value: "",
			Usage: "Slack WebHook url",
		},
	}
}
