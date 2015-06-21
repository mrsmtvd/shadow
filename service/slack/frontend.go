package slack

import (
	"github.com/GeertJohan/go.rice"
	"github.com/kihamo/shadow/service/frontend"
)

func (s *SlackService) GetTemplateBox() *rice.Box {
	return rice.MustFindBox("../slack/templates")
}

func (s *SlackService) GetFrontendMenu() *frontend.FrontendMenu {
	return &frontend.FrontendMenu{
		Name: "Slack",
		Url:  "/slack",
	}
}

func (s *SlackService) SetFrontendHandlers(router *frontend.Router) {
	router.GET(s, "/slack", &IndexHandler{})
}
