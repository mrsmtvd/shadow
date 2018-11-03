package http

import (
	"net/http"

	"github.com/kihamo/shadow/components/logging"
	misc "github.com/kihamo/shadow/misc/http"
)

func ServerMiddleware(log logging.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response := misc.NewResponse(w)

			next.ServeHTTP(response, r)

			statusCode := response.StatusCode()
			if statusCode == 0 {
				statusCode = http.StatusOK
			}

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

			if statusCode/100 == 5 {
				log.Error(http.StatusText(statusCode), fields)
			} else {
				log.Info(http.StatusText(statusCode), fields)
			}
		})
	}
}
