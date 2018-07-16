package internal

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/alexedwards/scs"
	"github.com/justinas/alice"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/dashboard/auth"
	"github.com/kihamo/shadow/components/logger"
	"github.com/kihamo/shadow/components/metrics"
)

const (
	DefaultClientName  = "unknown"
	DefaultHandlerName = "unknown"
)

func ContextMiddleware(application shadow.Application, router *Router, config config.Component, logger logger.Logger, renderer *Renderer, sessionManager *scs.Manager) alice.Constructor {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			writer := dashboard.NewResponse(w)
			request := dashboard.NewRequest(r)
			session := NewSession(sessionManager.Load(r), w)
			route := dashboard.RouteFromContext(r.Context())

			ctx := context.WithValue(r.Context(), dashboard.ApplicationContextKey, application)
			ctx = context.WithValue(ctx, dashboard.ConfigContextKey, config)
			ctx = context.WithValue(ctx, dashboard.LoggerContextKey, logger)
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

func LoggerMiddleware() alice.Constructor {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)

			statusCode := dashboard.ResponseFromContext(r.Context()).StatusCode()

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

			response := dashboard.ResponseFromContext(r.Context())

			handlerName := DefaultHandlerName
			route := dashboard.RouteFromContext(r.Context())
			if route != nil {
				handlerName = fmt.Sprintf("%s/%s", route.(*RouteItem).Component().Name(), route.HandlerName())
			}

			status := metrics.StatusOK
			if response.StatusCode()/100 >= 4 { // 4xx + 5xx
				status = metrics.StatusError
			}

			metrics.MetricRequestsTotal.With(
				"handler", handlerName,
				"protocol", metrics.ProtocolHTTP,
				"client_name", DefaultClientName,
			).Inc()

			metrics.MetricRequestSizeBytes.With(
				"handler", handlerName,
				"protocol", metrics.ProtocolHTTP,
				"client_name", DefaultClientName,
			).Add(float64(r.ContentLength))

			metrics.MetricResponseTimeSeconds.With(
				"handler", handlerName,
				"protocol", metrics.ProtocolHTTP,
				"client_name", DefaultClientName,
				"status", status,
			).UpdateSince(now)

			metrics.MetricResponseSizeBytes.With(
				"handler", handlerName,
				"protocol", metrics.ProtocolHTTP,
				"client_name", DefaultClientName,
				"status", status,
			).Add(float64(response.Length()))
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
