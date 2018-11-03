package internal

import (
	"context"
	"net/http"

	"github.com/alexedwards/scs"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/dashboard/auth"
)

func ContextMiddleware(application shadow.Application, router *Router, config config.Component, renderer *Renderer, sessionManager *scs.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			writer := dashboard.NewResponse(w)
			request := dashboard.NewRequest(r)
			session := NewSession(sessionManager.Load(r), w)
			route := dashboard.RouteFromContext(r.Context())

			ctx := context.WithValue(r.Context(), dashboard.ApplicationContextKey, application)
			ctx = context.WithValue(ctx, dashboard.ConfigContextKey, config)
			ctx = context.WithValue(ctx, dashboard.RenderContextKey, renderer)
			ctx = context.WithValue(ctx, dashboard.ResponseContextKey, writer)
			ctx = context.WithValue(ctx, dashboard.RouterContextKey, router)
			ctx = context.WithValue(ctx, dashboard.SessionContextKey, session)
			ctx = context.WithValue(ctx, dashboard.RequestContextKey, request)

			if route != nil {
				if routeItem, ok := route.(*RouteItem); ok {
					ctx = context.WithValue(ctx, dashboard.ComponentContextKey, routeItem.Component())
				}
			}

			request = request.WithContext(ctx)
			next.ServeHTTP(writer, request.Original())
		})
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

			if len(auth.GetProviders()) > 0 && !request.User().IsAuthorized() {
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
