package dashboard

import (
	"context"
	"net/http"
	"time"

	"github.com/justinas/alice"
)

type ResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *ResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *ResponseWriter) Write(data []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}

	return w.ResponseWriter.Write(data)
}

func (w *ResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *ResponseWriter) GetStatusCode() int {
	return w.status
}

func ContextMiddleware(c *Component) alice.Constructor {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), ConfigContextKey, c.config)
			ctx = context.WithValue(ctx, LoggerContextKey, c.logger)
			ctx = context.WithValue(ctx, RenderContextKey, c.renderer)
			ctx = context.WithValue(ctx, RequestContextKey, r)
			ctx = context.WithValue(ctx, ResponseContextKey, w)

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func BasicAuthMiddleware(c *Component) alice.Constructor {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.mutex.RLock()
			token := c.authToken
			c.mutex.RUnlock()

			if token == "" || r.Header.Get("Authorization") == token {
				next.ServeHTTP(w, r)
				return
			}

			w.Header().Set("WWW-Authenticate", `Basic realm="Shadow Security Zone"`)
			w.WriteHeader(401)
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

			switch writer.GetStatusCode() {
			case 500:
				logger.Error("Internal error", fields)

			case 404:
				logger.Warn("Not found", fields)

			default:
				logger.Info("Request", fields)
			}
		})
	}
}

func MetricsMiddleware(_ *Component) alice.Constructor {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			now := time.Now()

			next.ServeHTTP(w, r)

			if metricHandlerExecuteTime != nil {
				metricHandlerExecuteTime.ObserveDurationByTime(now)
			}
		})
	}
}
