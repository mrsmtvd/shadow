package dashboard

import (
	"context"
	"fmt"
	"net/http"
	"runtime"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

type Router struct {
	httprouter.Router

	defaultChain alice.Chain
	authChain    alice.Chain
}

func NewRouter(c *Component) *Router {
	r := &Router{}
	r.RedirectTrailingSlash = true
	r.RedirectFixedPath = true
	r.HandleMethodNotAllowed = true

	// chains
	r.defaultChain = alice.New(
		ContextMiddleware(c),
		MetricsMiddleware(c),
		LoggerMiddleware(c),
	)
	r.authChain = r.defaultChain.Append(BasicAuthMiddleware(c))

	return r
}

func (r *Router) setMiddleware(h http.Handler) http.Handler {
	var chain alice.Chain

	if authHandler, ok := h.(HandlerAuth); ok && authHandler.IsAuth() {
		chain = r.authChain
	} else {
		chain = r.defaultChain
	}

	return chain.Then(h)
}

func (r *Router) SetPanicHandler(h http.Handler) {
	r.PanicHandler = func(pw http.ResponseWriter, pr *http.Request, pe interface{}) {
		stack := make([]byte, 4096)
		stack = stack[:runtime.Stack(stack, false)]
		_, file, line, _ := runtime.Caller(6)

		r.setMiddleware(http.HandlerFunc(func(hw http.ResponseWriter, hr *http.Request) {
			ctx := context.WithValue(hr.Context(), PanicContextKey, &PanicError{
				error: pe,
				stack: string(stack),
				file:  file,
				line:  line,
			})

			hr = hr.WithContext(ctx)
			h.ServeHTTP(hw, hr)
		})).ServeHTTP(pw, pr)
	}
}

func (r *Router) SetNotAllowedHandler(h http.Handler) {
	r.MethodNotAllowed = r.setMiddleware(h)
}

func (r *Router) SetNotFoundHandler(h http.Handler) {
	r.NotFound = r.setMiddleware(h)
}

func (r *Router) HandlerFunc(m, p string, h http.HandlerFunc) {
	r.Router.Handler(m, p, r.authChain.ThenFunc(h))
}

func (r *Router) Handle(m, p string, h interface{}) {
	if h1, ok := h.(http.Handler); ok {
		r.Router.Handler(m, p, r.authChain.Then(h1))
	} else if h2, ok := h.(http.HandlerFunc); ok {
		r.Router.Handler(m, p, r.authChain.ThenFunc(h2))
	} else if h3, ok := h.(http.FileSystem); ok {
		r.Router.ServeFiles(p, h3)
	} else {
		panic(fmt.Sprintf("Unknown handler type %s %s %T", m, p, h))
	}
}
