package dashboard

import (
	"net/http"
)

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
				if route.ComponentName == "" {
					route.ComponentName = component.GetName()
				}

				router.AddRoute(route)
			}
		}
	}

	mainRouter := &Route{
		ComponentName: c.GetName(),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/"+c.GetName(), http.StatusMovedPermanently)
		}),
		Path:   "/",
		Direct: true,
	}

	router.AddRoute(mainRouter)

	return router, nil
}
