package dashboard

import (
	"net/http"
)

type Route struct {
	Methods []string
	Path    string
	Direct  bool
	Handler interface{}
	Auth    bool
}

type hasRoute interface {
	GetDashboardRoutes() []*Route
}

func (c *Component) getServeMux() (http.Handler, error) {
	// init routes
	components, err := c.application.GetComponents()
	if err != nil {
		return nil, err
	}

	router := NewRouter(c)

	router.SetPanicHandler(&PanicHandler{})
	router.SetForbiddenHandler(&ForbiddenHandler{})
	router.SetNotFoundHandler(&NotFoundHandler{})
	router.SetNotAllowedHandler(&MethodNotAllowedHandler{})

	for _, component := range components {
		if componentRoute, ok := component.(hasRoute); ok {
			for _, route := range componentRoute.GetDashboardRoutes() {
				path := route.Path
				if !route.Direct {
					path = "/" + component.GetName() + path
				}

				for _, method := range route.Methods {
					router.Handle(method, path, route.Handler, route.Auth)
				}
			}
		}
	}

	mainHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/"+c.GetName(), http.StatusMovedPermanently)
	})

	router.Handle("*", "/", mainHandler, false)

	// init session
	mux := NewSessionManager()(router)

	return mux, nil
}
