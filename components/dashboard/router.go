package dashboard

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"sync"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/kihamo/shadow/components/logger"
)

var httpMethods = []string{
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

	mutex  sync.RWMutex
	chain  alice.Chain
	logger logger.Logger
	routes []*Route
}

type Route struct {
	ComponentName string
	HandlerName   string
	Handler       interface{}
	Methods       []string
	Path          string
	Direct        bool
	Auth          bool
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

		r.chain.Then(http.HandlerFunc(func(hw http.ResponseWriter, hr *http.Request) {
			ctx := context.WithValue(hr.Context(), PanicContextKey, &PanicError{
				error: pe,
				stack: string(stack),
				file:  file,
				line:  line,
			})

			hr = hr.WithContext(ctx)
			panicHandler.ServeHTTP(hw, hr)
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

func (r *Router) GetRoutes() []Route {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	routes := make([]Route, 0, len(r.routes))
	for _, route := range r.routes {
		routes = append(routes, *route)
	}

	return routes
}

func (r *Router) AddMiddleware(m alice.Constructor) {
	r.chain = r.chain.Append(m)
}

func (r *Router) AddRoute(route *Route) {
	if !route.Direct {
		route.Path = "/" + route.ComponentName + route.Path
	}

	if len(route.Methods) == 0 {
		route.Methods = httpMethods
	}

	if route.HandlerName == "" {
		t := reflect.TypeOf(route.Handler)

		if t.Kind() == reflect.Ptr {
			route.HandlerName = t.Elem().Name()
		} else {
			route.HandlerName = t.Name()
		}
	}

	var handler http.Handler

	if h0, ok := route.Handler.(RouterHandler); ok {
		handler = FromRouteHandler(h0)
	} else if h1, ok := route.Handler.(http.Handler); ok {
		handler = h1
	} else if h2, ok := route.Handler.(http.HandlerFunc); ok {
		handler = h2
	} else if h3, ok := route.Handler.(http.FileSystem); ok {
		r.Router.ServeFiles(route.Path, h3)

		// TODO: set auth, metrics
		return
	} else {
		panic(fmt.Sprintf("Unknown handler type %s.%s for path %s", route.ComponentName, route.HandlerName, route.Path))
	}

	for i, method := range route.Methods {
		r.logger.Debug("Add handler", map[string]interface{}{
			"component": route.ComponentName,
			"handler":   route.HandlerName,
			"method":    method,
			"path":      route.Path,
			"direct":    route.Direct,
			"auth":      route.Auth,
		})

		localChan := alice.New(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := context.WithValue(r.Context(), RouteContextKey, route)
				next.ServeHTTP(w, r.WithContext(ctx))
			})
		})

		r.Router.Handler(method, route.Path, localChan.Extend(r.chain).Then(handler))

		if i == 0 {
			r.mutex.Lock()
			r.routes = append(r.routes, route)
			r.mutex.Unlock()
		}
	}
}
