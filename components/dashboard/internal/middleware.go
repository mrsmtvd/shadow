package internal

import (
	"context"
	"net/http"
	"time"

	"github.com/alexedwards/scs"
	"github.com/justinas/alice"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/dashboard/auth"
	"github.com/kihamo/shadow/components/logger"
)

func ContextMiddleware(router *Router, config config.Component, logger logger.Logger, renderer *Renderer, sessionManager *scs.Manager) alice.Constructor {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			writer := dashboard.NewResponse(w)
			request := dashboard.NewRequest(r)
			session := NewSession(sessionManager.Load(r), w)

			ctx := context.WithValue(r.Context(), dashboard.ConfigContextKey, config)
			ctx = context.WithValue(ctx, dashboard.LoggerContextKey, logger)
			ctx = context.WithValue(ctx, dashboard.RenderContextKey, renderer)
			ctx = context.WithValue(ctx, dashboard.ResponseContextKey, writer)
			ctx = context.WithValue(ctx, dashboard.RouterContextKey, router)
			ctx = context.WithValue(ctx, dashboard.SessionContextKey, session)
			ctx = context.WithValue(ctx, dashboard.RequestContextKey, request)

			request = request.WithContext(ctx)
			next.ServeHTTP(writer, request.Original())
		})
	}
}

func LoggerMiddleware() alice.Constructor {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)

			statusCode := dashboard.ResponseFromContext(r.Context()).GetStatusCode()

			fields := map[string]interface{}{
				"remote-addr":    r.RemoteAddr,
				"method":         r.Method,
				"request-uri":    r.RequestURI,
				"prote":          r.Proto,
				"code":           statusCode,
				"content-length": r.ContentLength,
				"referer":        r.Referer(),
				"user-agent":     r.UserAgent(),
			}

			log := dashboard.LoggerFromContext(r.Context())

			if statusCode/100 == 5 {
				log.Error(http.StatusText(statusCode), fields)
			} else {
				log.Info(http.StatusText(statusCode), fields)
			}
		})
	}
}

func MetricsMiddleware() alice.Constructor {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			now := time.Now()

			next.ServeHTTP(w, r)

			route := dashboard.RouteFromContext(r.Context())
			if route != nil {
				metricHandlerExecuteTime.With(
					"source", route.(*RouteItem).Source(),
					"handler", route.HandlerName(),
				).UpdateSince(now)
			}
		})
	}
}

func AuthorizationMiddleware() alice.Constructor {
	return func(next http.Handler) http.Handler {
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
}
