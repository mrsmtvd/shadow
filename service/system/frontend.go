package system

import (
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/resource/logger"
	"github.com/kihamo/shadow/resource/mail"
	"github.com/kihamo/shadow/resource/workers"
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
		{
			Name: "Configuration",
			Url:  "/system/config",
		},
		{
			Name: "Environment",
			Url:  "/system/environment",
		},
	}

	if s.application.HasResource("workers") {
		menu = append(menu, &frontend.FrontendMenu{
			Name: "Workers",
			Url:  "/system/workers",
		})
	}

	if s.application.HasResource("mail") {
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
	handlerConfig := &ConfigHandler{
		config: s.config,
		logger: logger.NewOrNop(s.GetName(), s.application),
	}
	router.GET(s, "/system/config", handlerConfig)
	router.POST(s, "/system/config", handlerConfig)

	router.GET(s, "/system/environment", &EnvironmentHandler{})

	if resourceWorkers, err := s.application.GetResource("workers"); err == nil {
		router.GET(s, "/system/workers", &WorkersHandler{
			workers: resourceWorkers.(*workers.Resource),
		})
	}

	if resourceMail, err := s.application.GetResource("mail"); err == nil {
		handlerMail := &MailHandler{
			mail: resourceMail.(*mail.Resource),
		}

		router.GET(s, "/system/mail", handlerMail)
		router.POST(s, "/system/mail", handlerMail)
	}
}
