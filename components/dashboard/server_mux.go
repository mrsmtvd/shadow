package dashboard

import (
	"net/http"

	"github.com/alexedwards/scs/engine/memstore"
	"github.com/alexedwards/scs/session"
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
					router.Handle(method, path, route.Handler)
				}
			}
		}
	}

	mainHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/"+c.GetName(), http.StatusMovedPermanently)
	})

	router.Handle(http.MethodGet, "/", mainHandler)
	router.Handle(http.MethodHead, "/", mainHandler)
	router.Handle(http.MethodPost, "/", mainHandler)
	router.Handle(http.MethodPut, "/", mainHandler)
	router.Handle(http.MethodPatch, "/", mainHandler)
	router.Handle(http.MethodDelete, "/", mainHandler)
	router.Handle(http.MethodConnect, "/", mainHandler)
	router.Handle(http.MethodOptions, "/", mainHandler)
	router.Handle(http.MethodTrace, "/", mainHandler)

	// init session
	session.CookieName = "shadow.token"

	sessionEngine := memstore.New(0)
	sessionManager := session.Manage(sessionEngine)

	mux := sessionManager(router)

	return mux, nil
}
