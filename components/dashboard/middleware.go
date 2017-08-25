package dashboard

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/justinas/alice"
)

func ContextMiddleware(c *Component) alice.Constructor {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), ConfigContextKey, c.config)
			ctx = context.WithValue(ctx, LoggerContextKey, c.logger)
			ctx = context.WithValue(ctx, RenderContextKey, c.renderer)
			ctx = context.WithValue(ctx, RequestContextKey, r)
			ctx = context.WithValue(ctx, ResponseContextKey, w)
			ctx = context.WithValue(ctx, RouterContextKey, c.router)

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func BasicAuthMiddleware(c *Component) alice.Constructor {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			config := ConfigFromContext(r.Context())

			checkUsername := config.GetString(ConfigAuthUser)
			if checkUsername == "" {
				next.ServeHTTP(w, r)
				return
			}

			checkPassword := config.GetString(ConfigAuthPassword)

			username, password, ok := r.BasicAuth()
			if ok && checkUsername == username && checkPassword == password {
				next.ServeHTTP(w, r)
				return
			}

			w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s Security Zone"`, c.application.GetName()))
			w.WriteHeader(http.StatusUnauthorized)
		})
	}
}

func LoggerMiddleware(c *Component) alice.Constructor {
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
