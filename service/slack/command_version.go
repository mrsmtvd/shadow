package slack

import (
	"github.com/nlopes/slack"
)

type VersionCommand struct {
	AbstractSlackCommand
}

func (c *VersionCommand) GetName() string {
	return "version"
}

func (c *VersionCommand) GetDescription() string {
	return "Возвращает текущую версию бота"
}

func (c *VersionCommand) Run(m *slack.MessageEvent, args ...string) {
	service := c.Service.(*SlackService)
	c.SendMessagef(m.Channel, "%s v.%s build %s / %s", service.Bot.Name, c.Application.Version, c.Application.Build, service.config.GetString("env"))
}
