package system

import (
	"strings"

	"github.com/kihamo/shadow/service/slack"
	sl "github.com/nlopes/slack"
)

type ServicesCommand struct {
	slack.AbstractSlackCommand
}

func (c *ServicesCommand) GetName() string {
	return "services"
}

func (c *ServicesCommand) GetDescription() string {
	return "Список подключенных служб"
}

func (c *ServicesCommand) Run(m *sl.MessageEvent, args ...string) {
	service := c.Service.(*SystemService)

	values := []string{}

	for _, s := range service.Application.GetServices() {
		values = append(values, s.GetName())
	}

	c.SendMessage(m.Channel, strings.Join(values, "\n"))
	return
}
