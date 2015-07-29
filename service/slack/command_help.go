package slack

import (
	"github.com/nlopes/slack"
)

type HelpCommand struct {
	AbstractSlackCommand
}

func (c *HelpCommand) GetName() string {
	return "help"
}

func (c *HelpCommand) GetDescription() string {
	return "Справка по всем командам"
}

func (c *HelpCommand) Run(m *slack.MessageEvent, args ...string) {
	params := slack.NewPostMessageParameters()
	params.Attachments = []slack.Attachment{slack.Attachment{
		Color:  "good",
		Fields: []slack.AttachmentField{},
	}}

	service := c.Service.(*SlackService)
	if len(args) == 0 {
		for name := range service.Commands {
			params.Attachments[0].Fields = append(params.Attachments[0].Fields, slack.AttachmentField{
				Title: service.Commands[name].GetName(),
				Value: service.Commands[name].GetDescription(),
			})
		}
	} else if command, ok := service.Commands[args[0]]; ok {
		params.Attachments[0].Fields = append(params.Attachments[0].Fields, slack.AttachmentField{
			Title: command.GetName(),
			Value: command.GetDescription(),
		})
	} else {
		c.SendMessagef(m.Channel, "Команда %s не найдена", args[0])
		return
	}

	c.SendPostMessage(m.Channel, "Мои команды", params)
}
