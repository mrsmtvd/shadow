package slack

import (
	"github.com/nlopes/slack"
)

type HelloCommand struct {
	AbstractSlackCommand
}

func (c *HelloCommand) GetName() string {
	return "hello"
}

func (c *HelloCommand) GetDescription() string {
	return "Приветствие. Так я реагирую на любое упоминание моего имени."
}

func (c *HelloCommand) Run(m *slack.MessageEvent, args ...string) {
	c.SendMessage(m.Channel, "Привет! Спроси у меня")
}
