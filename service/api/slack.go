package api

import (
	"github.com/kihamo/shadow/service/slack"
)

func (s *ApiService) GetSlackCommands() []slack.SlackCommand {
	return []slack.SlackCommand{
		&ApiCommand{},
	}
}
