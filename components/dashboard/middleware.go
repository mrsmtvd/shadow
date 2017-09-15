package dashboard

import (
	"context"
	"net/http"
	"time"

	"github.com/alexedwards/scs"
	"github.com/justinas/alice"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/logger"
)

func ContextMiddleware(router *Router, config *config.Component, logger logger.Logger, renderer *Renderer, sessionManager *scs.Manager) alice.Constructor {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			writer := NewResponse(w)
			request := NewRequest(r)
			session := NewSession(sessionManager.Load(r), w)

			ctx := context.WithValue(r.Context(), ConfigContextKey, config)
			ctx = context.WithValue(ctx, LoggerContextKey, logger)
			ctx = context.WithValue(ctx, RenderContextKey, renderer)
			ctx = context.WithValue(ctx, ResponseContextKey, writer)
			ctx = context.WithValue(ctx, RouterContextKey, router)
			ctx = context.WithValue(ctx, SessionContextKey, session)
			ctx = context.WithValue(ctx, RequestContextKey, request)
			r = r.WithContext(ctx)

			// TODO: dirty hack
			request.original = r

			next.ServeHTTP(writer, r)
		})
	}
}

func LoggerMiddleware() alice.Constructor {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)

			statusCode := ResponseFromContext(r.Context()).GetStatusCode()

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

			logger := LoggerFromContext(r.Context())

			if statusCode/100 == 5 {
				logger.Error(http.StatusText(statusCode), fields)
			} else {
				logger.Info(http.StatusText(statusCode), fields)
			}
		})
	}
}

func MetricsMiddleware() alice.Constructor {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			now := time.Now()

			next.ServeHTTP(w, r)

			route := RouteFromContext(r.Context())
			if route != nil {
				metricHandlerExecuteTime.With(
					"component", route.ComponentName,
					"handler", route.HandlerName,
				).UpdateSince(now)
			}
		})
	}
}

func AuthorizationMiddleware() alice.Constructor {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			route := RouteFromContext(r.Context())
			if route != nil && route.Auth {
				session := SessionFromContext(r.Context())
				if session == nil {
					panic("Session isn't set in context")
				}

				auth, err := session.GetString(SessionUsername)
				if err != nil || auth == "" {
					request := RequestFromContext(r.Context())
					session := SessionFromContext(r.Context())
					if request != nil && session != nil && !request.IsAjax() && request.IsGet() {
						session.PutString(SessionLastURL, request.URL().Path)
					}

					http.Redirect(w, r, "/dashboard/login", http.StatusFound)
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
