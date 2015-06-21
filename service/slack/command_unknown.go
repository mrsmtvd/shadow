package slack

import (
	"github.com/nlopes/slack"
)

type UnknownCommand struct {
	AbstractSlackCommand
}

func (c *UnknownCommand) GetName() string {
	return "unknown"
}

func (c *UnknownCommand) GetDescription() string {
	return "Реакция на неизвестную команду"
}

func (c *UnknownCommand) Run(m *slack.MessageEvent, args ...string) {
	c.SendMessagef(m.ChannelId, "Я не знаю такой команды :wink:. Попробуй набрать: <@%s>: help", c.Service.(*SlackService).Bot.Name)
}
