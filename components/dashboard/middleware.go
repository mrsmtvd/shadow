package dashboard

import (
	"context"
	"net/http"
	"time"

	"github.com/justinas/alice"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/logger"
)

func ContextMiddleware(router *Router, config *config.Component, logger logger.Logger, renderer *Renderer) alice.Constructor {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), ConfigContextKey, config)
			ctx = context.WithValue(ctx, LoggerContextKey, logger)
			ctx = context.WithValue(ctx, RenderContextKey, renderer)
			ctx = context.WithValue(ctx, RequestContextKey, r)
			ctx = context.WithValue(ctx, ResponseContextKey, w)
			ctx = context.WithValue(ctx, RouterContextKey, router)

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func LoggerMiddleware() alice.Constructor {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			writer := &ResponseWriter{
				ResponseWriter: w,
			}

			next.ServeHTTP(writer, r)

			fields := map[string]interface{}{
				"remote-addr":    r.RemoteAddr,
				"method":         r.Method,
				"request-uri":    r.RequestURI,
				"prote":          r.Proto,
				"code":           writer.GetStatusCode(),
				"content-length": r.ContentLength,
				"referer":        r.Referer(),
				"user-agent":     r.UserAgent(),
			}

			logger := LoggerFromContext(r.Context())

			if writer.GetStatusCode()/100 == 5 {
				logger.Error(http.StatusText(writer.GetStatusCode()), fields)
			} else {
				logger.Info(http.StatusText(writer.GetStatusCode()), fields)
			}
		})
	}
}

func MetricsMiddleware(c *Component) alice.Constructor {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			now := time.Now()

			next.ServeHTTP(w, r)

			if c.application.HasComponent("metrics") {
				metricHandlerExecuteTime.UpdateSince(now)
			}
		})
	}
}
