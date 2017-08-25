package dashboard

import (
	"context"
	"fmt"
	"net/http"
	"runtime"

	"github.com/alexedwards/scs/session"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/kihamo/shadow/components/logger"
)

type Router struct {
	httprouter.Router

	Forbidden http.Handler

	chain  alice.Chain
	logger logger.Logger
}

func NewRouter(c *Component) *Router {
	r := &Router{}
	r.RedirectTrailingSlash = true
	r.RedirectFixedPath = true
	r.HandleMethodNotAllowed = true

	// chains
	r.chain = alice.New(
		ContextMiddleware(r, c.config, c.logger, c.renderer),
		MetricsMiddleware(c),
		LoggerMiddleware(),
		r.authMiddleware,
	)

	r.logger = c.logger

	return r
}

func (r *Router) authMiddleware(next http.Handler) http.Handler {
	if r.Forbidden == nil {
		return next
	}

	if authHandler, ok := next.(HandlerAuth); !ok || !authHandler.IsAuth() {
		return next
	}

	if _, ok := next.(http.HandlerFunc); !ok && r.Forbidden == next {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, rq *http.Request) {
		auth, err := session.GetString(rq, "username")

		if err == nil && auth != "" {
			next.ServeHTTP(w, rq)
		} else {
			r.Forbidden.ServeHTTP(w, rq)
		}
	})
}

func (r *Router) SetPanicHandler(h http.Handler) {
	r.PanicHandler = func(pw http.ResponseWriter, pr *http.Request, pe interface{}) {
		stack := make([]byte, 4096)
		stack = stack[:runtime.Stack(stack, false)]
		_, file, line, _ := runtime.Caller(6)

		r.chain.Then(http.HandlerFunc(func(hw http.ResponseWriter, hr *http.Request) {
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

func (r *Router) SetForbiddenHandler(h http.Handler) {
	r.Forbidden = r.chain.Then(h)
}

func (r *Router) SetNotFoundHandler(h http.Handler) {
	r.NotFound = r.chain.Then(h)
}

func (r *Router) SetNotAllowedHandler(h http.Handler) {
	r.MethodNotAllowed = r.chain.Then(h)
}

func (r *Router) HandlerFunc(m, p string, h http.HandlerFunc) {
	r.Router.Handler(m, p, r.chain.ThenFunc(h))
}

func (r *Router) Handle(m, p string, h interface{}) {
	if h1, ok := h.(http.Handler); ok {
		r.Router.Handler(m, p, r.chain.Then(h1))
	} else if h2, ok := h.(http.HandlerFunc); ok {
		r.Router.Handler(m, p, r.chain.Then(h2))
	} else if h3, ok := h.(http.FileSystem); ok {
		r.Router.ServeFiles(p, h3)
	} else {
		panic(fmt.Sprintf("Unknown handler type %s %s %T", m, p, h))
	}

	if authHandler, ok := h.(HandlerAuth); ok && authHandler.IsAuth() {
		r.logger.Debugf("Add security handler for %s %s", m, p)
	} else {
		r.logger.Debugf("Add handler for %s %s", m, p)
	}
}
