package internal

import (
	"net/http"

	"github.com/kihamo/shadow/components/dashboard"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
)

func (c *Component) DashboardMiddleware() []func(http.Handler) http.Handler {
	return []func(http.Handler) http.Handler{
		c.ServerMiddleware,
	}
}

func (c *Component) ServerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tracer := c.Tracer()

		if tracer == nil {
			next.ServeHTTP(w, r)
		} else {
			nethttp.MiddlewareFunc(tracer, func(w http.ResponseWriter, r *http.Request) {
				if span := opentracing.SpanFromContext(r.Context()); span != nil {
					if sc, ok := span.Context().(jaeger.SpanContext); ok {
						w.Header().Set("trace-id", sc.TraceID().String())
					}
				}

				next.ServeHTTP(w, r)
			}, nethttp.OperationNameFunc(func(r *http.Request) string {
				var handlerName string

				route := dashboard.RouteFromContext(r.Context())
				if route != nil {
					handlerName = route.HandlerName()
				}

				return "HTTP " + r.Method + ": " + handlerName
			})).ServeHTTP(w, r)
		}
	})
}
