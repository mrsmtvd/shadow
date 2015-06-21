package system

import (
	"fmt"
	"strings"

	"github.com/kihamo/shadow/service/slack"
	sl "github.com/nlopes/slack"
)

type ConfigCommand struct {
	slack.AbstractSlackCommand
}

func (c *ConfigCommand) GetName() string {
	return "config"
}

func (c *ConfigCommand) GetDescription() string {
	return "Информация о конфигурации системы"
}

func (c *ConfigCommand) Run(m *sl.MessageEvent, args ...string) {
	service := c.Service.(*SystemService)

	if len(args) == 0 {
		values := []string{}

		for name := range service.Config.GetAll() {
			values = append(values, fmt.Sprintf("*%s* = %s", name, fmt.Sprint(service.Config.Get(name))))
		}

		c.SendMessage(m.ChannelId, strings.Join(values, "\n"))
		return
	}

	if !service.Config.Has(args[0]) {
		c.SendMessagef(m.ChannelId, "Настройки *%s* не существует", args[0])
		return
	}

	c.SendMessagef(m.ChannelId, "*%s* = %s", args[0], fmt.Sprint(service.Config.Get(args[0])))
}
