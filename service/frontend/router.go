package frontend

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/kihamo/shadow"
)

type Router struct {
	httprouter.Router
	application *shadow.Application
}

func NewRouter(application *shadow.Application) *Router {
	r := &Router{}
	r.RedirectTrailingSlash = true
	r.RedirectFixedPath = true
	r.HandleMethodNotAllowed = true
	r.application = application

	return r
}

func (r *Router) GET(s shadow.Service, path string, h interface{}) {
	r.Handle(s, "GET", path, h)
}

func (r *Router) POST(s shadow.Service, path string, h interface{}) {
	r.Handle(s, "POST", path, h)
}

func (r *Router) Handle(s shadow.Service, m, p string, h interface{}) {
	if h1, ok := h.(Handler); ok {
		h1.Init(r.application, s)

		r.Router.Handle(m, p, func(out http.ResponseWriter, in *http.Request, _ httprouter.Params) {
			out.Header().Set("Content-Type", "text/html; charset=utf-8")

			h1.InitRequest(out, in)
			h1.Handle()
			h1.Render()
		})
	} else if h2, ok := h.(http.Handler); ok {
		r.Router.Handler(m, p, h2)
	} else if h3, ok := h.(http.HandlerFunc); ok {
		r.Router.HandlerFunc(m, p, h3)
	} else if h4, ok := h.(httprouter.Handle); ok {
		r.Router.Handle(m, p, h4)
	} else {
		panic(fmt.Sprintf("Unknown handler type %s %s", m, p))
	}
}
