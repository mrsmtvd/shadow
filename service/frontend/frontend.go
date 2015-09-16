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
		Icon: "dashboard",
	}
}

func (s *FrontendService) SetFrontendHandlers(router *Router) {
	router.ServeFiles("/vendor/*filepath", &assetfs.AssetFS{
		Asset:    Asset,
		AssetDir: AssetDir,
		Prefix:   "public/vendor",
	})

	asset := &assetfs.AssetFS{
		Asset:    Asset,
		AssetDir: AssetDir,
		Prefix:   "public",
	}
	router.GET(s, "/favicon.svg", http.HandlerFunc(func(out http.ResponseWriter, in *http.Request) {
		http.FileServer(asset).ServeHTTP(out, in)
	}))
	router.GET(s, "/frontend.css", http.HandlerFunc(func(out http.ResponseWriter, in *http.Request) {
		http.FileServer(asset).ServeHTTP(out, in)
	}))

	router.GET(s, "/alerts", &AlertsHandler{})
	router.GET(s, "/", &IndexHandler{})
}
