package frontend

import (
	"net/http"

	"github.com/GeertJohan/go.rice"
)

func (s *FrontendService) GetTemplateBox() *rice.Box {
	return rice.MustFindBox("../frontend/templates")
}

func (s *FrontendService) GetFrontendMenu() *FrontendMenu {
	return &FrontendMenu{
		Name: "Main",
		Url:  "/",
	}
}

func (s *FrontendService) SetFrontendHandlers(router *Router) {
	box := rice.MustFindBox("public/css")
	router.ServeFiles("/css/*filepath", box.HTTPBox())

	box = rice.MustFindBox("public/fonts")
	router.ServeFiles("/fonts/*filepath", box.HTTPBox())

	box = rice.MustFindBox("public/js")
	router.ServeFiles("/js/*filepath", box.HTTPBox())

	router.GET(s, "/favicon.svg", http.HandlerFunc(func(out http.ResponseWriter, in *http.Request) {
		out.Write(s.boxStatic.MustBytes("favicon.svg"))
	}))

	router.GET(s, "/", &IndexHandler{})
}
