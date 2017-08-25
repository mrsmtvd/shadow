package dashboard

import (
	"context"
	originalHttp "net/http"
	"time"

	"github.com/justinas/alice"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard/http"
	"github.com/kihamo/shadow/components/logger"
)

func ContextMiddleware(router *Router, config *config.Component, logger logger.Logger, renderer *Renderer) alice.Constructor {
	return func(next originalHttp.Handler) originalHttp.Handler {
		return originalHttp.HandlerFunc(func(w originalHttp.ResponseWriter, r *originalHttp.Request) {
			writer := http.NewResponse(w)
			reader := http.NewRequest(r)

			ctx := context.WithValue(r.Context(), ConfigContextKey, config)
			ctx = context.WithValue(ctx, LoggerContextKey, logger)
			ctx = context.WithValue(ctx, RenderContextKey, renderer)
			ctx = context.WithValue(ctx, RequestContextKey, reader)
			ctx = context.WithValue(ctx, ResponseContextKey, writer)
			ctx = context.WithValue(ctx, RouterContextKey, router)

			r = r.WithContext(ctx)
			next.ServeHTTP(writer, r)
		})
	}
}

func LoggerMiddleware() alice.Constructor {
	return func(next originalHttp.Handler) originalHttp.Handler {
		return originalHttp.HandlerFunc(func(w originalHttp.ResponseWriter, r *originalHttp.Request) {
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
				logger.Error(originalHttp.StatusText(statusCode), fields)
			} else {
				logger.Info(originalHttp.StatusText(statusCode), fields)
			}
		})
	}
}

func MetricsMiddleware(c *Component) alice.Constructor {
	return func(next originalHttp.Handler) originalHttp.Handler {
		return originalHttp.HandlerFunc(func(w originalHttp.ResponseWriter, r *originalHttp.Request) {
			now := time.Now()

			next.ServeHTTP(w, r)

			if c.application.HasComponent("metrics") {
				metricHandlerExecuteTime.UpdateSince(now)
			}
		})
	}
}
