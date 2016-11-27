package frontend

import (
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
)

func (s *FrontendService) GetTemplates() *assetfs.AssetFS {
	return &assetfs.AssetFS{
		Asset:     Asset,
		AssetDir:  AssetDir,
		AssetInfo: AssetInfo,
		Prefix:    "templates",
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
		Asset:     Asset,
		AssetDir:  AssetDir,
		AssetInfo: AssetInfo,
		Prefix:    "public/vendor",
	})

	asset := &assetfs.AssetFS{
		Asset:     Asset,
		AssetDir:  AssetDir,
		AssetInfo: AssetInfo,
		Prefix:    "public",
	}
	router.GET(s, "/favicon.svg", http.HandlerFunc(func(out http.ResponseWriter, in *http.Request) {
		http.FileServer(asset).ServeHTTP(out, in)
	}))
	router.GET(s, "/frontend.css", http.HandlerFunc(func(out http.ResponseWriter, in *http.Request) {
		http.FileServer(asset).ServeHTTP(out, in)
	}))
	router.GET(s, "/frontend.js", http.HandlerFunc(func(out http.ResponseWriter, in *http.Request) {
		http.FileServer(asset).ServeHTTP(out, in)
	}))

	if s.application.HasResource("alerts") {
		router.GET(s, "/alerts", &AlertsHandler{})
	}

	router.GET(s, "/", &IndexHandler{})
}
