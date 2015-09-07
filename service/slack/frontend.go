package slack

import (
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/service/frontend"
)

func (s *SlackService) GetTemplates() *assetfs.AssetFS {
	return &assetfs.AssetFS{
		Asset:    Asset,
		AssetDir: AssetDir,
		Prefix:   "templates",
	}
}

func (s *SlackService) GetFrontendMenu() *frontend.FrontendMenu {
	return &frontend.FrontendMenu{
		Name: "Slack",
		Url:  "/slack",
		Icon: "slack",
	}
}

func (s *SlackService) SetFrontendHandlers(router *frontend.Router) {
	router.GET(s, "/slack", &IndexHandler{})
}
