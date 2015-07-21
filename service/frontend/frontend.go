package frontend

import (
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
)

func (s *FrontendService) GetTemplates() *assetfs.AssetFS {
	return &assetfs.AssetFS{
		Asset:    Asset,
		AssetDir: AssetDir,
		Prefix:   "templates",
	}
}

func (s *FrontendService) GetFrontendMenu() *FrontendMenu {
	return &FrontendMenu{
		Name: "Main",
		Url:  "/",
	}
}

func (s *FrontendService) SetFrontendHandlers(router *Router) {
	router.ServeFiles("/css/*filepath", &assetfs.AssetFS{
		Asset:    Asset,
		AssetDir: AssetDir,
		Prefix:   "public/css",
	})

	router.ServeFiles("/fonts/*filepath", &assetfs.AssetFS{
		Asset:    Asset,
		AssetDir: AssetDir,
		Prefix:   "public/fonts",
	})

	router.ServeFiles("/js/frontend/*filepath", &assetfs.AssetFS{
		Asset:    Asset,
		AssetDir: AssetDir,
		Prefix:   "public/js",
	})

	asset := &assetfs.AssetFS{
		Asset:    Asset,
		AssetDir: AssetDir,
		Prefix:   "public",
	}
	router.GET(s, "/favicon.svg", http.HandlerFunc(func(out http.ResponseWriter, in *http.Request) {
		http.FileServer(asset).ServeHTTP(out, in)
	}))

	router.GET(s, "/", &IndexHandler{})
}
