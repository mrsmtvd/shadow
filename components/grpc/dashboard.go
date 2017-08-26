package grpc

import (
	"net/http"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/components/dashboard"
)

func (c *Component) GetTemplates() *assetfs.AssetFS {
	return &assetfs.AssetFS{
		Asset:     Asset,
		AssetDir:  AssetDir,
		AssetInfo: AssetInfo,
		Prefix:    "templates",
	}
}

func (c *Component) GetDashboardMenu() *dashboard.Menu {
	return &dashboard.Menu{
		Name: "gRPC",
		Url:  "/",
		Icon: "exchange",
	}
}

func (c *Component) GetDashboardRoutes() []*dashboard.Route {
	return []*dashboard.Route{
		{
			Methods: []string{http.MethodGet},
			Path:    "/assets/*filepath",
			Handler: &assetfs.AssetFS{
				Asset:     Asset,
				AssetDir:  AssetDir,
				AssetInfo: AssetInfo,
				Prefix:    "assets",
			},
		},
		{
			Methods: []string{http.MethodGet, http.MethodPost},
			Path:    "/",
			Handler: NewIndexHandler(c),
			Auth:    true,
		},
	}
}
