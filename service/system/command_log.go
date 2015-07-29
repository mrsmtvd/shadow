package system

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
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
	return "Show logs"
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

		c.SendPostMessage(m.Channel, "Available components", params)
		return
	}

	showItemsCount := MaxItems
	if len(args) == 2 {
		var err error
		if showItemsCount, err = strconv.Atoi(args[1]); err != nil {
			c.SendMessage(m.Channel, "Specified number of records to display is not a number")
			return
		}
	}

	component, ok := loggers[args[0]]
	if !ok {
		c.SendMessagef(m.Channel, "Component *%s* does't exists", args[0])
		return
	}

	if len(component) == 0 {
		c.SendMessagef(m.Channel, "Log empty for *%s* component", args[0])
		return
	}

	skip := -1
	if showItemsCount < len(component) {
		skip = len(component) - showItemsCount
	} else {
		showItemsCount = len(component)
	}

	var (
		color      string
		index      int64
		attachment sl.Attachment
	)

	params := sl.NewPostMessageParameters()
	params.Attachments = make([]sl.Attachment, showItemsCount)

	for i, c := range component {
		if i < skip {
			continue
		}

		switch c.Level {
		case logrus.DebugLevel:
			color = "#999999"
		case logrus.InfoLevel:
			color = "#5bc0de"
		case logrus.WarnLevel:
			color = "warning"
		case logrus.ErrorLevel:
			color = "danger"
		case logrus.FatalLevel:
			color = "danger"
		case logrus.PanicLevel:
			color = "danger"
		default:
			color = "good"
		}

		attachment = sl.Attachment{
			Color:  color,
			Title:  c.Time.Format(time.RFC1123Z),
			Text:   c.Message,
			Fields: []sl.AttachmentField{},
		}

		for field, value := range component[i].Data {
			if field == "component" {
				continue
			}

			attachment.Fields = append(attachment.Fields, sl.AttachmentField{
				Title: field,
				Value: fmt.Sprintf("%v", value),
				Short: true,
			})
		}

		params.Attachments[index] = attachment
		index = index + 1
	}

	c.SendPostMessage(m.Channel, fmt.Sprintf("Show last *%d* entries for *%s* component", showItemsCount, args[0]), params)
}
