package dashboard

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

type Router struct {
	httprouter.Router

	defaultChain alice.Chain
	authChain    alice.Chain
}

func NewRouter(component *Component) *Router {
	r := &Router{}
	r.RedirectTrailingSlash = true
	r.RedirectFixedPath = true
	r.HandleMethodNotAllowed = true

	// chains
	r.defaultChain = alice.New(
		MetricsMiddleware(component),
		LoggerMiddleware(component),
	)
	r.authChain = r.defaultChain.Append(BasicAuthMiddleware(component))

	return r
}

func (r *Router) getInitHandler(h Handler) http.Handler {
	var chain alice.Chain

	if authHandler, ok := h.(HandlerAuth); ok && authHandler.IsAuth() {
		chain = r.authChain
	} else {
		chain = r.defaultChain
	}

	return chain.ThenFunc(func(response http.ResponseWriter, request *http.Request) {
		h.SetRequest(request)
		h.SetResponse(response)
		h.Handle()

		if templateHandler, ok := h.(HandlerTemplate); ok {
			templateHandler.Render()
		}
	})
}

func (r *Router) SetPanicHandler(h Handler) {
	var chain alice.Chain

	if auth, ok := h.(HandlerAuth); ok && auth.IsAuth() {
		chain = r.authChain
	} else {
		chain = r.defaultChain
	}

	r.PanicHandler = func(w http.ResponseWriter, r *http.Request, error interface{}) {
		chain.ThenFunc(func(response http.ResponseWriter, request *http.Request) {
			h.SetRequest(request)
			h.SetResponse(response)

			if panicHandler, ok := h.(HandlerPanic); ok {
				panicHandler.SetError(error)
			}

			h.Handle()

			if templateHandler, ok := h.(HandlerTemplate); ok {
				templateHandler.Render()
			}
		}).ServeHTTP(w, r)
	}
}

func (r *Router) SetNotAllowedHandler(h Handler) {
	r.MethodNotAllowed = r.getInitHandler(h)
}

func (r *Router) SetNotFoundHandler(h Handler) {
	r.NotFound = r.getInitHandler(h)
}

func (r *Router) HandlerFunc(m, p string, h http.HandlerFunc) {
	r.Router.Handler(m, p, r.authChain.ThenFunc(h))
}

func (r *Router) Handle(m, p string, h interface{}) {
	if h1, ok := h.(Handler); ok {
		r.Router.Handler(m, p, r.getInitHandler(h1))
	} else if h2, ok := h.(http.Handler); ok {
		r.Router.Handler(m, p, r.authChain.Then(h2))
	} else if h3, ok := h.(http.HandlerFunc); ok {
		r.Router.Handler(m, p, r.authChain.ThenFunc(h3))
	} else if h4, ok := h.(http.FileSystem); ok {
		r.Router.ServeFiles(p, h4)
	} else {
		panic(fmt.Sprintf("Unknown handler type %s %s %T", m, p, h))
	}
}
