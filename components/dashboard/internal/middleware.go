package internal

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/mrsmtvd/shadow/components/dashboard"
)

func ContextMiddleware(router dashboard.Router, renderer dashboard.Renderer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			writer := dashboard.NewResponse(w)
			request := dashboard.NewRequest(r)
			route := dashboard.RouteFromContext(r.Context())

			ctx := dashboard.ContextWithRender(r.Context(), renderer)
			ctx = dashboard.ContextWithResponse(ctx, writer)
			ctx = dashboard.ContextWithRouter(ctx, router)
			ctx = dashboard.ContextWithRequest(ctx, request)

			if route != nil {
				if routeItem, ok := route.(*RouteItem); ok {
					ctx = dashboard.ContextWithTemplateNamespace(ctx, routeItem.Component().Name())
				}
			}

			request = request.WithContext(ctx)
			next.ServeHTTP(writer, request.Original())
		})
	}
}

func SessionMiddleware(sessionManager *scs.SessionManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return sessionManager.LoadAndSave(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session := NewSession(sessionManager, r)
			ctx := dashboard.ContextWithSession(r.Context(), session)

			next.ServeHTTP(w, r.WithContext(ctx))

			session.Flush()
		}))
	}
}

func AuthorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := dashboard.RouteFromContext(r.Context())
		if route != nil && route.Auth() {
			request := dashboard.RequestFromContext(r.Context())
			if request == nil {
				panic("Request isn't set in context")
			}

			if !request.User().IsAuthorized() {
				if !request.IsAjax() && request.IsGet() {
					request.Session().PutString(dashboard.SessionLastURL, request.URL().Path)
				}

				http.Redirect(w, r, dashboard.AuthPath, http.StatusFound)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
