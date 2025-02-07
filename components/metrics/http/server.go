package http

import (
	"net/http"
	"time"

	"github.com/mrsmtvd/shadow/components/metrics"
	misc "github.com/mrsmtvd/shadow/misc/http"
)

const (
	defaultClientName  = "unknown"
	defaultHandlerName = "unknown"
)

func ServerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := misc.NewResponse(w)

		now := time.Now()

		// client name
		clientName := defaultClientName
		if r.Header.Get("User-Agent") != "" {
			clientName = "browser"
		}

		// TODO: handler name
		handlerName := defaultHandlerName

		metrics.MetricRequestsTotal.With(
			"handler", handlerName,
			"protocol", metrics.ProtocolHTTP,
			"client_name", clientName,
		).Inc()

		metrics.MetricRequestSizeBytes.With(
			"handler", handlerName,
			"protocol", metrics.ProtocolHTTP,
			"client_name", clientName,
		).Add(float64(r.ContentLength))

		next.ServeHTTP(response, r)

		status := metrics.StatusOK
		if response.StatusCode()/100 >= 4 { // 4xx + 5xx
			status = metrics.StatusError
		}

		metrics.MetricResponseTimeSeconds.With(
			"handler", handlerName,
			"protocol", metrics.ProtocolHTTP,
			"client_name", clientName,
			"status", status,
		).UpdateSince(now)

		metrics.MetricResponseSizeBytes.With(
			"handler", handlerName,
			"protocol", metrics.ProtocolHTTP,
			"client_name", clientName,
			"status", status,
		).Add(float64(response.Length()))
	})
}
