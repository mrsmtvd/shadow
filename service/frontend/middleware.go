package frontend

import (
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

func BasicAuthMiddleware(service *FrontendService) alice.Constructor {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			service.mutex.RLock()
			token := service.authToken
			service.mutex.RUnlock()

			if token == "" || r.Header.Get("Authorization") == token {
				next.ServeHTTP(w, r)
				return
			}

			w.Header().Set("WWW-Authenticate", `Basic realm="Shadow Security Zone"`)
			w.WriteHeader(401)
		})
	}
}

func LoggerMiddleware(service *FrontendService) alice.Constructor {
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

			switch writer.GetStatusCode() {
			case 500:
				service.logger.Error("Internal error", fields)

			case 404:
				service.logger.Warn("Not found", fields)

			default:
				service.logger.Info("Request", fields)
			}
		})
	}
}

func MetricsMiddleware(service *FrontendService) alice.Constructor {
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
