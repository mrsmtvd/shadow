package api

import (
	slacks "github.com/kihamo/shadow-slack/service"
)

func (s *ApiService) GetSlackCommands() []slacks.SlackCommand {
	return []slacks.SlackCommand{
		&ApiCommand{},
	}
}
