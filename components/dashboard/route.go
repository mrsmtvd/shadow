package dashboard

import (
	"net/http"
	"reflect"
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

type Route interface {
	HandlerName() string
	Handler() interface{}
	Methods() []string
	Path() string
	Auth() bool
}

type HasRoutes interface {
	DashboardRoutes() []Route
}

type RouteSimple struct {
	handlerName string
	handler     interface{}
	methods     []string
	path        string
	auth        bool
}

func NewRoute(path string, handler interface{}) *RouteSimple {
	return &RouteSimple{
		handler: handler,
		methods: httpMethods,
		path:    path,
	}
}

func (r *RouteSimple) HandlerName() string {
	if r.handlerName == "" {
		t := reflect.TypeOf(r.handler)

		if t.Kind() == reflect.Ptr {
			r.handlerName = t.Elem().Name()
		} else {
			r.handlerName = t.Name()
		}
	}

	return r.handlerName
}

func (r *RouteSimple) WithHandlerName(handlerName string) *RouteSimple {
	r.handlerName = handlerName
	return r
}

func (r *RouteSimple) Handler() interface{} {
	return r.handler
}

func (r *RouteSimple) WithHandler(handler interface{}) *RouteSimple {
	r.handler = handler
	return r
}

func (r *RouteSimple) Methods() []string {
	if len(r.methods) == 0 {
		return httpMethods
	}

	return r.methods
}

func (r *RouteSimple) WithMethods(methods []string) *RouteSimple {
	r.methods = methods
	return r
}

func (r *RouteSimple) Path() string {
	return r.path
}

func (r *RouteSimple) WithPath(path string) *RouteSimple {
	r.path = path
	return r
}

func (r *RouteSimple) Auth() bool {
	return r.auth
}

func (r *RouteSimple) WithAuth(auth bool) *RouteSimple {
	r.auth = auth
	return r
}
