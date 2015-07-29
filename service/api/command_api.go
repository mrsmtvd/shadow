package api

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"

	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/service/slack"
	sl "github.com/nlopes/slack"
	"gopkg.in/jcelliott/turnpike.v2"
)

type ApiCommand struct {
	slack.AbstractSlackCommand
	Service *ApiService

	client *turnpike.Client
}

func (c *ApiCommand) GetName() string {
	return "api"
}

func (c *ApiCommand) GetDescription() string {
	return "Вызов процедур Api"
}

func (c *ApiCommand) Init(s shadow.Service, a *shadow.Application) {
	c.AbstractSlackCommand.Init(s, a)
	c.Service = s.(*ApiService)
}

func (c *ApiCommand) Run(m *sl.MessageEvent, args ...string) {
	if len(args) == 0 {
		params := sl.NewPostMessageParameters()
		params.Attachments = []sl.Attachment{sl.Attachment{
			Color:  "good",
			Fields: []sl.AttachmentField{},
		}}

		for _, procedure := range c.Service.GetProcedures() {
			params.Attachments[0].Fields = append(params.Attachments[0].Fields, sl.AttachmentField{
				Title: procedure.GetName(),
			})
		}

		c.SendPostMessage(m.Channel, "Available api procedures", params)
		return
	}

	if !c.Service.HasProcedure(args[0]) {
		c.SendMessagef(m.Channel, "Procedure *%s* does't exists", args[0])
		return
	}

	var (
		result *turnpike.Result
		err    error
		out    bytes.Buffer
	)

	apiArgs := make(procedureArgs, 0)
	apiKwargs := make(procedureKwargs, 0)

	apiSet := flag.NewFlagSet("", flag.ContinueOnError)
	apiSet.Var(&apiArgs, "a", "Args")
	apiSet.Var(&apiKwargs, "k", "Kwargs")

	if err = apiSet.Parse(args[1:]); err != nil {
		c.SendMessage(m.Channel, err.Error())
		return
	}

	params := sl.NewPostMessageParameters()

	if c.client == nil {
		c.client, err = c.Service.GetClient()
	}

	if err == nil {
		result, err = c.client.Call(args[0], apiArgs, apiKwargs)
	}

	if err != nil {
		params.Attachments = []sl.Attachment{sl.Attachment{
			Color: "danger",
			Fields: []sl.AttachmentField{
				sl.AttachmentField{
					Title: fmt.Sprintf("Procedure %s call error", args[0]),
					Value: err.Error(),
				},
			},
		}}
	} else {
		r, _ := json.Marshal(result)
		json.Indent(&out, r, "", "\t")

		params.Attachments = []sl.Attachment{sl.Attachment{
			Color: "good",
			Title: fmt.Sprintf("Procedure %s call success", args[0]),
			Text:  out.String(),
		}}
	}

	c.SendPostMessage(m.Channel, "", params)
}
