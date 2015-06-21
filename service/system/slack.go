package system

import (
	"github.com/kihamo/shadow/service/slack"
)

func (s *SystemService) GetSlackCommands() []slack.SlackCommand {
	return []slack.SlackCommand{
		&ConfigCommand{},
		&LogCommand{},
	}
}
