package dashboard

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"runtime"

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

	Forbidden http.Handler

	chain  alice.Chain
	logger logger.Logger
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
	r.chain = alice.New(
		ContextMiddleware(r, c.config, c.logger, c.renderer, c.session),
		MetricsMiddleware(c),
		LoggerMiddleware(),
	)

	r.logger = c.logger

	return r
}

func (r *Router) getHandlerName(h interface{}) string {
	t := reflect.TypeOf(h)

	if t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	}

	return t.Name()
}

func (r *Router) setMetaMiddleware(hhandler http.Handler, route *Route) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, rq *http.Request) {
		ctx := context.WithValue(rq.Context(), RouteContextKey, route)

		hhandler.ServeHTTP(w, rq.WithContext(ctx))
	})
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

func (r *Router) AddRoute(route *Route) {
	if !route.Direct {
		route.Path = "/" + route.ComponentName + route.Path
	}

	if len(route.Methods) == 0 {
		route.Methods = httpMethods
	}

	if route.HandlerName == "" {
		route.HandlerName = r.getHandlerName(route.Handler)
	}

	var handler http.Handler

	if h0, ok := route.Handler.(RouterHandler); ok {
		handler = r.setAuthMiddleware(FromRouteHandler(h0), route.Auth)
	} else if h1, ok := route.Handler.(http.Handler); ok {
		handler = r.setAuthMiddleware(h1, route.Auth)
	} else if h2, ok := route.Handler.(http.HandlerFunc); ok {
		handler = r.setAuthMiddleware(h2, route.Auth)
	} else if h3, ok := route.Handler.(http.FileSystem); ok {
		r.Router.ServeFiles(route.Path, h3)

		// TODO: set auth, metrics
		return
	} else {
		panic(fmt.Sprintf("Unknown handler type %s.%s for path %s", route.ComponentName, route.HandlerName, route.Path))
	}

	handler = r.chain.Then(handler)

	for _, method := range route.Methods {
		r.logger.Debug("Add handler", map[string]interface{}{
			"component": route.ComponentName,
			"handler":   route.HandlerName,
			"method":    method,
			"path":      route.Path,
			"direct":    route.Direct,
			"auth":      route.Auth,
		})

		r.Router.Handler(method, route.Path, r.setMetaMiddleware(handler, route))
	}
}
