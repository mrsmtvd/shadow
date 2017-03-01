package dashboard

import (
	"net/http"
)

type Route struct {
	Methods []string
	Path    string
	Direct  bool
	Handler interface{}
}

type hasRoute interface {
	GetDashboardRoutes() []*Route
}

func (c *Component) loadRoutes() error {
	components, err := c.application.GetComponents()
	if err != nil {
		return err
	}

	http.DefaultServeMux = http.NewServeMux()

	c.router = NewRouter(c)

	c.router.SetPanicHandler(&PanicHandler{})
	c.router.SetNotAllowedHandler(&MethodNotAllowedHandler{})
	c.router.SetNotFoundHandler(&NotFoundHandler{})

	for _, component := range components {
		if componentRoute, ok := component.(hasRoute); ok {
			for _, route := range componentRoute.GetDashboardRoutes() {
				path := route.Path
				if !route.Direct {
					path = "/" + component.GetName() + path
				}

				for _, method := range route.Methods {
					c.router.Handle(method, path, route.Handler)
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

	return nil
}
