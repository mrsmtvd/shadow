package internal

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"sync"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/logger"
)

type Router struct {
	httprouter.Router

	mutex  sync.RWMutex
	chain  alice.Chain
	logger logger.Logger
	routes []dashboard.Route
}

type RouterHandler interface {
	ServeHTTP(*dashboard.Response, *dashboard.Request)
}

func FromRouteHandler(h RouterHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var wq *dashboard.Response

		if resp, ok := w.(*dashboard.Response); ok {
			wq = resp
		} else {
			wq = dashboard.NewResponse(w)
		}

		h.ServeHTTP(wq, dashboard.NewRequest(r))
	})
}

func NewRouter(c *Component) *Router {
	r := &Router{}
	r.RedirectTrailingSlash = true
	r.RedirectFixedPath = true
	r.HandleMethodNotAllowed = true

	// chains
	r.chain = alice.New()
	r.logger = c.logger

	return r
}

func (r *Router) SetPanicHandler(h RouterHandler) {
	panicHandler := FromRouteHandler(h)

	r.PanicHandler = func(pw http.ResponseWriter, pr *http.Request, pe interface{}) {
		stack := make([]byte, 4096)
		stack = stack[:runtime.Stack(stack, false)]
		_, file, line, _ := runtime.Caller(6)

		fmt.Println(dashboard.SessionFromContext(pr.Context()))

		r.chain.Then(http.HandlerFunc(func(hw http.ResponseWriter, hr *http.Request) {
			ctx := context.WithValue(hr.Context(), dashboard.PanicContextKey, &dashboard.PanicError{
				Error: pe,
				Stack: string(stack),
				File:  file,
				Line:  line,
			})

			panicHandler.ServeHTTP(hw, hr.WithContext(ctx))
		})).ServeHTTP(pw, pr)
	}
}

func (r *Router) SetNotFoundHandler(h RouterHandler) {
	handler := FromRouteHandler(h)

	r.NotFound = http.HandlerFunc(func(w http.ResponseWriter, rq *http.Request) {
		r.chain.Then(handler).ServeHTTP(w, rq)
	})
}

func (r *Router) SetNotAllowedHandler(h RouterHandler) {
	handler := FromRouteHandler(h)

	r.MethodNotAllowed = http.HandlerFunc(func(w http.ResponseWriter, rq *http.Request) {
		r.chain.Then(handler).ServeHTTP(w, rq)
	})
}

func (r *Router) GetRoutes() []dashboard.Route {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	routes := make([]dashboard.Route, 0, len(r.routes))
	for _, route := range r.routes {
		routes = append(routes, route)
	}

	return routes
}

func (r *Router) AddMiddleware(m alice.Constructor) {
	r.chain = r.chain.Append(m)
}

func (r *Router) AddRoute(route dashboard.Route) {
	var handler http.Handler

	if h0, ok := route.Handler().(RouterHandler); ok {
		handler = FromRouteHandler(h0)
	} else if h1, ok := route.Handler().(http.Handler); ok {
		handler = h1
	} else if h2, ok := route.Handler().(http.HandlerFunc); ok {
		handler = h2
	} else if h3, ok := route.Handler().(http.FileSystem); ok {
		r.Router.ServeFiles(route.Path(), h3)

		r.mutex.Lock()
		r.routes = append(r.routes, route)
		r.mutex.Unlock()

		// TODO: set auth, metrics
		return
	} else {
		panic(fmt.Sprintf("Unknown handler type %s.%s for path %s", route.ComponentName(), route.HandlerName(), route.Path()))
	}

	for i, method := range route.Methods() {
		r.logger.Debug("Add handler", map[string]interface{}{
			"component": route.ComponentName,
			"handler":   route.HandlerName,
			"method":    method,
			"path":      route.Path,
			"auth":      route.Auth,
		})

		localChan := alice.New(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := context.WithValue(r.Context(), dashboard.RouteContextKey, route)
				next.ServeHTTP(w, r.WithContext(ctx))
			})
		})

		r.Router.Handle(method, route.Path(), func(w http.ResponseWriter, rq *http.Request, p httprouter.Params) {
			values := rq.URL.Query()
			for _, param := range p {
				values.Add(":"+param.Key, param.Value)
			}
			rq.URL.RawQuery = values.Encode()

			localChan.Extend(r.chain).Then(handler).ServeHTTP(w, rq)
		})

		if i == 0 {
			r.mutex.Lock()
			r.routes = append(r.routes, route)
			r.mutex.Unlock()
		}
	}
}

func (r *Router) NotFoundServeHTTP(w http.ResponseWriter, rq *http.Request) {
	r.NotFound.ServeHTTP(w, rq)
}

func (r *Router) MethodNotAllowedServeHTTP(w http.ResponseWriter, rq *http.Request) {
	r.MethodNotAllowed.ServeHTTP(w, rq)
}
