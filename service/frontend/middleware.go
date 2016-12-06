package frontend

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/justinas/alice"
	"github.com/kihamo/shadow/resource/metrics"
	"github.com/rs/xlog"
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
	user := service.config.GetString("frontend.auth-user")
	if user == "" {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
			})
		}
	}

	password := service.config.GetString("frontend.auth-password")
	token := "Basic " + base64.StdEncoding.EncodeToString([]byte(user+":"+password))

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Authorization") == token {
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

			message := fmt.Sprintf("%s \"%s %s %s\" %d %d \"%s\" \"%s\"", r.RemoteAddr, r.Method, r.RequestURI, r.Proto, writer.GetStatusCode(), r.ContentLength, r.Referer(), r.UserAgent())

			fields := xlog.F{
				"method":      r.Method,
				"request-uri": r.RequestURI,
				"code":        writer.GetStatusCode(),
			}

			switch writer.GetStatusCode() {
			case 500:
				service.Logger.Error(message, fields)

			case 404:
				service.Logger.Warn(message, fields)

			default:
				service.Logger.Info(message, fields)
			}
		})
	}
}

func MetricsMiddleware(service *FrontendService) alice.Constructor {
	resourceMetrics, err := service.application.GetResource("metrics")
	if err == nil {
		timer := resourceMetrics.(*metrics.Metrics).NewTimer(MetricHandlerExecuteTime)

		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				beforeTime := time.Now()
				next.ServeHTTP(w, r)
				timer.UpdateSince(beforeTime)
			})
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
}
