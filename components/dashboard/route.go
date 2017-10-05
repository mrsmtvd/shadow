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
	ComponentName() string
	HandlerName() string
	Handler() interface{}
	Methods() []string
	Path() string
	Auth() bool
}

type HasRoutes interface {
	GetDashboardRoutes() []Route
}

type RouteItem struct {
	componentName string
	handlerName   string
	handler       interface{}
	methods       []string
	path          string
	auth          bool
}

func NewRoute(componentName string, methods []string, path string, handler interface{}, handlerName string, auth bool) Route {
	return RouteItem{
		componentName: componentName,
		handlerName:   handlerName,
		handler:       handler,
		methods:       methods,
		path:          path,
		auth:          auth,
	}
}

func (r RouteItem) ComponentName() string {
	return r.componentName
}

func (r RouteItem) HandlerName() string {
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

func (r RouteItem) Handler() interface{} {
	return r.handler
}

func (r RouteItem) Methods() []string {
	if len(r.methods) == 0 {
		return httpMethods
	}

	return r.methods
}

func (r RouteItem) Path() string {
	return r.path
}

func (r RouteItem) Auth() bool {
	return r.auth
}
