package internal

import (
	"net/http"
	"runtime"
	"runtime/debug"
	"sync"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/mrsmtvd/shadow/components/dashboard"
	"github.com/mrsmtvd/shadow/components/logging"
)

const (
	DefaultCallerSkip = 4
)

type Router struct {
	httprouter.Router

	mutex      sync.RWMutex
	chain      alice.Chain
	logger     logging.Logger
	routes     []dashboard.Route
	callerSkip int
}

type RouterHandler interface {
	ServeHTTP(*dashboard.Response, *dashboard.Request)
}

type RouterMixHandler interface {
	ServeHTTP(http.ResponseWriter, *dashboard.Request)
}

func FromRouteHandler(h RouterHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(dashboard.NewResponse(w), dashboard.NewRequest(r))
	})
}

func NewRouter(l logging.Logger, skip int) *Router {
	r := &Router{
		chain:      alice.New(),
		logger:     l,
		callerSkip: skip,
	}
	r.RedirectTrailingSlash = true
	r.RedirectFixedPath = true
	r.HandleMethodNotAllowed = true

	return r
}

func (r *Router) SetPanicHandlerCallerSkip(skip int) {
	r.mutex.Lock()
	r.callerSkip = skip
	r.mutex.Unlock()
}

// TODO: двойной вызов мидлварей для хендлера происходит
func (r *Router) SetPanicHandler(h RouterHandler) {
	panicHandler := FromRouteHandler(h)

	r.PanicHandler = func(pw http.ResponseWriter, pr *http.Request, pe interface{}) {
		r.mutex.RLock()
		skip := r.callerSkip
		r.mutex.RUnlock()

		_, file, line, _ := runtime.Caller(skip)

		r.chain.Then(http.HandlerFunc(func(hw http.ResponseWriter, hr *http.Request) {
			panicError := &dashboard.PanicError{
				Error: pe,
				Stack: debug.Stack(),
				File:  file,
				Line:  line,
			}

			ctx := dashboard.ContextWithPanic(hr.Context(), panicError)

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

func (r *Router) Routes() []dashboard.Route {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	routes := make([]dashboard.Route, len(r.routes))
	copy(routes, r.routes)

	return routes
}

func (r *Router) addMiddleware(m func(next http.Handler) http.Handler) {
	r.chain = r.chain.Append(alice.Constructor(m))
}

func (r *Router) addRoute(route dashboard.Route) {
	var handler http.Handler

	switch h := route.Handler().(type) {
	case RouterHandler:
		handler = FromRouteHandler(h)
	case RouterMixHandler:
		handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, dashboard.NewRequest(r))
		})
	case http.HandlerFunc:
		handler = h
	case func(http.ResponseWriter, *http.Request):
		handler = http.HandlerFunc(h)
	case http.Handler:
		handler = h
	case http.FileSystem:
		r.Router.ServeFiles(route.Path(), h)

		r.mutex.Lock()
		r.routes = append(r.routes, route)
		r.mutex.Unlock()

		return
	default:
		panic("Unknown handler type " + route.HandlerName() + " for path " + route.Path())
	}

	for i, method := range route.Methods() {
		r.logger.Debug("Add handler",
			"handler", route.HandlerName(),
			"method", method,
			"path", route.Path(),
			"auth", route.Auth(),
		)

		localChan := alice.New(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				r = r.WithContext(dashboard.ContextWithRoute(r.Context(), route))
				next.ServeHTTP(w, r)
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

func (r *Router) InternalErrorServeHTTP(w http.ResponseWriter, rq *http.Request, e error) {
	r.PanicHandler(w, rq, e)
}
