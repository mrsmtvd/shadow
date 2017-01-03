package dashboard

import (
	"net/http"
)

type Route struct {
	Methods []string
	Path    string
	Handler interface{}
}

type hasRoute interface {
	GetDashboardRoutes() []*Route
}

func (c *Component) loadRoutes() {
	http.DefaultServeMux = http.NewServeMux()

	c.router = NewRouter(c)

	panicHandler := &PanicHandler{
		logger: c.logger,
	}
	panicHandler.SetRenderer(c.renderer)
	c.router.SetPanicHandler(panicHandler)

	methodNotAllowedHandler := &MethodNotAllowedHandler{}
	methodNotAllowedHandler.SetRenderer(c.renderer)
	c.router.SetNotAllowedHandler(methodNotAllowedHandler)

	notFoundHandler := &NotFoundHandler{}
	notFoundHandler.SetRenderer(c.renderer)
	c.router.SetNotFoundHandler(notFoundHandler)

	for _, component := range c.application.GetComponents() {
		if componentRoute, ok := component.(hasRoute); ok {
			for _, route := range componentRoute.GetDashboardRoutes() {
				if templateRoute, ok := route.Handler.(HandlerTemplate); ok {
					templateRoute.SetRenderer(c.renderer)
				}

				for _, method := range route.Methods {
					c.router.Handle(method, "/"+component.GetName()+route.Path, route.Handler)
				}
			}
		}
	}

	mainHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/"+c.GetName(), http.StatusMovedPermanently)
	})

	c.router.Handle(http.MethodGet, "/", mainHandler)
	c.router.Handle(http.MethodHead, "/", mainHandler)
	c.router.Handle(http.MethodPost, "/", mainHandler)
	c.router.Handle(http.MethodPut, "/", mainHandler)
	c.router.Handle(http.MethodPatch, "/", mainHandler)
	c.router.Handle(http.MethodDelete, "/", mainHandler)
	c.router.Handle(http.MethodConnect, "/", mainHandler)
	c.router.Handle(http.MethodOptions, "/", mainHandler)
	c.router.Handle(http.MethodTrace, "/", mainHandler)
}
