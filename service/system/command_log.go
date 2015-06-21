package system

import (
	"fmt"
	"strings"

	"strconv"
	"time"

	"github.com/kihamo/shadow/service/slack"
	sl "github.com/nlopes/slack"
)

type LogCommand struct {
	slack.AbstractSlackCommand
}

func (c *LogCommand) GetName() string {
	return "log"
}

func (c *LogCommand) GetDescription() string {
	return "Логи системы"
}

func (c *LogCommand) Run(m *sl.MessageEvent, args ...string) {
	loggers := loggerHook.GetRecords()

	if len(args) == 0 {
		params := sl.NewPostMessageParameters()
		params.Attachments = []sl.Attachment{sl.Attachment{
			Color:  "good",
			Fields: []sl.AttachmentField{},
		}}

		for name := range loggers {
			params.Attachments[0].Fields = append(params.Attachments[0].Fields, sl.AttachmentField{
				Title: name,
			})
		}

		c.SendPostMessage(m.ChannelId, "Доступные компоненты", params)
		return
	}

	showItemsCount := MaxItems
	if len(args) == 2 {
		var err error
		if showItemsCount, err = strconv.Atoi(args[1]); err != nil {
			c.SendMessage(m.ChannelId, "Значение количества выводимых записей не является числом")
			return
		}
	}

	component, ok := loggers[args[0]]
	if !ok {
		c.SendMessagef(m.ChannelId, "Компонента *%s* не существует", args[0])
		return
	}

	skip := -1
	if showItemsCount < len(component) {
		skip = len(component) - showItemsCount
	}

	values := []string{}
	for i := range component {
		if i < skip {
			continue
		}

		values = append(values, fmt.Sprintf("%s [%s] %s", component[i].Time.Format(time.RFC1123Z), strings.ToUpper(component[i].Level.String()), component[i].Message))
	}

	c.SendMessage(m.ChannelId, strings.Join(values, "\n"))
}
