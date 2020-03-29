package internal

import (
	"net/http"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/mail/internal/handlers"
)

func (c *Component) DashboardTemplates() *assetfs.AssetFS {
	return dashboard.TemplatesFromAssetFS(c)
}

func (c *Component) DashboardMenu() dashboard.Menu {
	return dashboard.NewMenu("Mail").
		WithURL("/" + c.Name() + "/send/").
		WithIcon("envelope")
}

func (c *Component) DashboardRoutes() []dashboard.Route {
	return []dashboard.Route{
		dashboard.NewRoute("/"+c.Name()+"/send/", handlers.NewSendHandler(c)).
			WithMethods([]string{http.MethodGet, http.MethodPost}).
			WithAuth(true),
	}
}
