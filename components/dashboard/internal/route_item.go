package internal

import (
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/dashboard"
)

type RouteItem struct {
	dashboard.Route

	component shadow.Component
	route     dashboard.Route
}

func NewRouteItem(route dashboard.Route, component shadow.Component) *RouteItem {
	return &RouteItem{
		route:     route,
		component: component,
	}
}

func (r *RouteItem) Component() shadow.Component {
	return r.component
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
