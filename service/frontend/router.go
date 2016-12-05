package frontend

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/kihamo/shadow"
)

type Router struct {
	httprouter.Router
	application *shadow.Application

	defaultChain alice.Chain
	authChain    alice.Chain
}

func NewRouter(service *FrontendService) *Router {
	r := &Router{}
	r.RedirectTrailingSlash = true
	r.RedirectFixedPath = true
	r.HandleMethodNotAllowed = true

	r.application = service.application

	// chains
	r.defaultChain = alice.New(
		MetricsMiddleware(service),
		LoggerMiddleware(service),
	)
	r.authChain = r.defaultChain.Append(BasicAuthMiddleware(service))

	return r
}

func (r *Router) getInitHandler(s shadow.Service, h Handler) http.Handler {
	var chain alice.Chain

	if authHandler, ok := h.(HandlerAuth); ok && authHandler.IsAuth() {
		chain = r.authChain
	} else {
		chain = r.defaultChain
	}

	h.Init(r.application, s)

	return chain.ThenFunc(func(out http.ResponseWriter, in *http.Request) {
		out.Header().Set("Content-Type", "text/html; charset=utf-8")

		h.InitRequest(out, in)
		h.Handle()
		h.Render()
	})
}

func (r *Router) SetPanicHandler(s shadow.Service, h Handler) {
	var chain alice.Chain

	if auth, ok := h.(HandlerAuth); ok && auth.IsAuth() {
		chain = r.authChain
	} else {
		chain = r.defaultChain
	}

	h.Init(r.application, s)

	r.PanicHandler = func(w http.ResponseWriter, r *http.Request, error interface{}) {
		chain.ThenFunc(func(out http.ResponseWriter, in *http.Request) {
			out.Header().Set("Content-Type", "text/html; charset=utf-8")
			h.InitRequest(out, in)

			if panicHandler, ok := h.(HandlerPanic); ok {
				panicHandler.SetError(error)
			}

			h.Handle()
			h.Render()
		}).ServeHTTP(w, r)
	}
}

func (r *Router) SetNotAllowedHandler(s shadow.Service, h Handler) {
	r.MethodNotAllowed = r.getInitHandler(s, h)
}

func (r *Router) SetNotFoundHandler(s shadow.Service, h Handler) {
	r.NotFound = r.getInitHandler(s, h)
}

func (r *Router) GET(s shadow.Service, path string, h interface{}) {
	r.Handle(s, "GET", path, h)
}

func (r *Router) POST(s shadow.Service, path string, h interface{}) {
	r.Handle(s, "POST", path, h)
}

func (r *Router) Handle(s shadow.Service, m, p string, h interface{}) {
	if h1, ok := h.(Handler); ok {
		r.Router.Handler(m, p, r.getInitHandler(s, h1))
	} else if h2, ok := h.(http.Handler); ok {
		r.Router.Handler(m, p, r.authChain.Then(h2))
	} else if h3, ok := h.(http.HandlerFunc); ok {
		r.Router.Handler(m, p, r.authChain.ThenFunc(h3))
	} else {
		panic(fmt.Sprintf("Unknown handler type %s %s", m, p))
	}
}
