package http

import (
	"net/http"

	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/tracing"
	misc "github.com/kihamo/shadow/misc/http"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
)

func ServerMiddleware(tracer opentracing.Tracer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response := misc.NewResponse(w)

			nethttp.MiddlewareFunc(tracer, func(w http.ResponseWriter, r *http.Request) {
				span := opentracing.SpanFromContext(r.Context())

				if span != nil {
					if sc, ok := span.Context().(jaeger.SpanContext); ok {
						w.Header().Set(tracing.ResponseTraceIDHeader, sc.TraceID().String())
					}
				}

				next.ServeHTTP(w, r)

				if span != nil {
					statusCode := response.StatusCode()
					if statusCode > 0 && statusCode/100 == 5 {
						ext.Error.Set(span, true)
					}
				}
			}, nethttp.OperationNameFunc(func(r *http.Request) string {
				var handlerName string

				route := dashboard.RouteFromContext(r.Context())
				if route != nil {
					handlerName = route.HandlerName()
				}

				return "HTTP " + r.Method + ": " + handlerName
			})).ServeHTTP(response, r)
		})
	}
}
