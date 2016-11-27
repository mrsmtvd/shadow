package system

import (
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/service/frontend"
)

func (s *SystemService) GetTemplates() *assetfs.AssetFS {
	return &assetfs.AssetFS{
		Asset:     Asset,
		AssetDir:  AssetDir,
		AssetInfo: AssetInfo,
		Prefix:    "templates",
	}
}

func (s *SystemService) GetFrontendMenu() *frontend.FrontendMenu {
	menu := []*frontend.FrontendMenu{
		&frontend.FrontendMenu{
			Name: "Environment",
			Url:  "/system/environment",
		},
		&frontend.FrontendMenu{
			Name: "Configuration",
			Url:  "/system/config",
		},
	}

	if s.Application.HasResource("logger") {
		menu = append(menu, &frontend.FrontendMenu{
			Name: "Logs",
			Url:  "/system/logs",
		})
	}

	if s.Application.HasResource("workers") {
		menu = append(menu, &frontend.FrontendMenu{
			Name: "Workers",
			Url:  "/system/workers",
		})
	}

	if s.Application.HasResource("mail") {
		menu = append(menu, &frontend.FrontendMenu{
			Name: "Mail",
			Url:  "/system/mail",
		})
	}

	return &frontend.FrontendMenu{
		Name:    "System",
		SubMenu: menu,
		Icon:    "gear",
	}
}

func (s *SystemService) SetFrontendHandlers(router *frontend.Router) {
	router.GET(s, "/system/config", &ConfigHandler{})
	router.GET(s, "/system/environment", &EnvironmentHandler{})

	if s.Application.HasResource("logger") {
		router.GET(s, "/system/logs", &LogsHandler{})
	}

	if s.Application.HasResource("workers") {
		router.GET(s, "/system/workers", &WorkersHandler{})
	}

	if s.Application.HasResource("mail") {
		handlerMail := &MailHandler{}

		router.GET(s, "/system/mail", handlerMail)
		router.POST(s, "/system/mail", handlerMail)
	}
}
