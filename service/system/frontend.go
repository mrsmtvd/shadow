package system

import (
	"github.com/GeertJohan/go.rice"
	"github.com/kihamo/shadow/service/frontend"
)

func (s *SystemService) GetTemplateBox() *rice.Box {
	return rice.MustFindBox("../system/templates")
}

func (s *SystemService) GetFrontendMenu() *frontend.FrontendMenu {
	return &frontend.FrontendMenu{
		Name: "System",
		SubMenu: []*frontend.FrontendMenu{
			&frontend.FrontendMenu{
				Name: "Config",
				Url:  "/system/config",
			},
			&frontend.FrontendMenu{
				Name: "Logs",
				Url:  "/system/logs",
			},
		},
	}
}

func (s *SystemService) SetFrontendHandlers(router *frontend.Router) {
	router.GET(s, "/system/config", &ConfigHandler{})

	if s.Application.HasResource("logger") {
		router.GET(s, "/system/logs", &LogsHandler{})
	}
}
