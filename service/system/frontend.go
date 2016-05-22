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

	if s.Application.HasResource("tasks") {
		menu = append(menu, &frontend.FrontendMenu{
			Name: "Tasks",
			Url:  "/system/tasks",
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

	if s.Application.HasResource("logger") {
		router.GET(s, "/system/logs", &LogsHandler{})
	}

	if s.Application.HasResource("tasks") {
		router.GET(s, "/system/tasks", &TasksHandler{})
	}

	if s.Application.HasResource("mail") {
		handlerMail := &MailHandler{}

		router.GET(s, "/system/mail", handlerMail)
		router.POST(s, "/system/mail", handlerMail)
	}
}
