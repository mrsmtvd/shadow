package aws

import (
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/service/frontend"
)

func (s *AwsService) GetTemplates() *assetfs.AssetFS {
	return &assetfs.AssetFS{
		Asset:    Asset,
		AssetDir: AssetDir,
		Prefix:   "templates",
	}
}

func (s *AwsService) GetFrontendMenu() *frontend.FrontendMenu {
	return &frontend.FrontendMenu{
		Name: "Aws",
		Url:  "/aws",
	}
}

func (s *AwsService) SetFrontendHandlers(router *frontend.Router) {
	router.GET(s, "/aws", &IndexHandler{})
}
