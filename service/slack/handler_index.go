package slack

import (
	"github.com/kihamo/shadow/service/frontend"
	"github.com/nlopes/slack"
)

type IndexHandler struct {
	frontend.AbstractFrontendHandler
}

func (h *IndexHandler) Handle() {
	service := h.Service.(*SlackService)

	h.SetTemplate("index.tpl.html")
	h.View.Context["PageTitle"] = "Slack"
	h.View.Context["Commands"] = service.Commands

	var (
		name     string
		channels []slack.Channel
	)

	if service.Bot != nil {
		name = service.Bot.Name
	}

	if service.Rtm != nil {
		channels, _ = service.Rtm.GetChannels(true)
	}

	h.View.Context["Slack"] = map[string]interface{}{
		"name":     name,
		"channels": channels,
	}
}
