package internal

import (
	"github.com/kihamo/shadow/components/dashboard"
)

type RouteItem struct {
	dashboard.Route

	source string
	route  dashboard.Route
}

func NewRouteItem(route dashboard.Route, source string) *RouteItem {
	if source == "" {
		source = "unknown"
	}

	return &RouteItem{
		route:  route,
		source: source,
	}
}

func (r *RouteItem) Source() string {
	return r.source
}

func (r *RouteItem) HandlerName() string {
	return r.route.HandlerName()
}

func (r *RouteItem) Handler() interface{} {
	return r.route.Handler()
}

func (r *RouteItem) Methods() []string {
	return r.route.Methods()
}

func (r *RouteItem) Path() string {
	return r.route.Path()
}

func (r *RouteItem) Auth() bool {
	return r.route.Auth()
}
