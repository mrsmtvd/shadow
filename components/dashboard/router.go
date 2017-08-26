package dashboard

import (
	"context"
	"fmt"
	"net/http"
	"runtime"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/kihamo/shadow/components/logger"
)

var httpMethods = [9]string{
	http.MethodGet,
	http.MethodHead,
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
	http.MethodConnect,
	http.MethodOptions,
	http.MethodTrace,
}

type Router struct {
	httprouter.Router

	Forbidden http.Handler

	chain  alice.Chain
	logger logger.Logger
}

type RouterHandler interface {
	ServeHTTP(*Response, *Request)
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
	)

	r.logger = c.logger

	return r
}

func (r *Router) setAuthMiddleware(h http.Handler, a bool) http.Handler {
	if !a || r.Forbidden == nil {
		return h
	}

	if _, ok := h.(http.HandlerFunc); !ok && r.Forbidden == h {
		return h
	}

	return http.HandlerFunc(func(w http.ResponseWriter, rq *http.Request) {
		auth, err := SessionFromContext(rq.Context()).GetString(SessionUsername)

		if err == nil && auth != "" {
			h.ServeHTTP(w, rq)
		} else {
			r.Forbidden.ServeHTTP(w, rq)
		}
	})
}

func (r *Router) SetPanicHandler(h RouterHandler) {
	panicHadler := FromRouteHandler(h)

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
			panicHadler.ServeHTTP(hw, hr)
		})).ServeHTTP(pw, pr)
	}
}

func (r *Router) SetForbiddenHandler(h RouterHandler) {
	r.Forbidden = r.chain.Then(FromRouteHandler(h))
}

func (r *Router) SetNotFoundHandler(h RouterHandler) {
	r.NotFound = r.chain.Then(FromRouteHandler(h))
}

func (r *Router) SetNotAllowedHandler(h RouterHandler) {
	r.MethodNotAllowed = r.chain.Then(FromRouteHandler(h))
}

func (r *Router) Handle(m, p string, h interface{}, a bool) {
	if m == "" || m == "*" {
		for _, method := range httpMethods {
			r.Handle(method, p, h, a)
		}

		return
	}

	var handler http.Handler

	if h0, ok := h.(RouterHandler); ok {
		handler = r.setAuthMiddleware(FromRouteHandler(h0), a)
	} else if h1, ok := h.(http.Handler); ok {
		handler = r.setAuthMiddleware(h1, a)
	} else if h2, ok := h.(http.HandlerFunc); ok {
		handler = r.setAuthMiddleware(h2, a)
	} else if h3, ok := h.(http.FileSystem); ok {
		r.Router.ServeFiles(p, h3)

		// TODO: set auth
		return
	} else {
		panic(fmt.Sprintf("Unknown handler type %s %s %T", m, p, h))
	}

	r.Router.Handler(m, p, r.chain.Then(handler))

	if a {
		r.logger.Debugf("Add security handler for %s %s", m, p)
	} else {
		r.logger.Debugf("Add handler for %s %s", m, p)
	}
}
