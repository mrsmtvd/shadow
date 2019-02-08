package http

import (
	"context"
	"net/http"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
)

type contextKey string

var (
	componentNameContextKey = contextKey("component")
	operationNameContextKey = contextKey("operation name")
	clientTraceContextKey   = contextKey("client trace")
)

func ComponentNameFromContext(ctx context.Context) string {
	v := ctx.Value(componentNameContextKey)
	if v != nil {
		return v.(string)
	}

	return ""
}

func ComponentNameToContext(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, componentNameContextKey, value)
}

func OperationNameFromContext(ctx context.Context) string {
	v := ctx.Value(operationNameContextKey)
	if v != nil {
		return v.(string)
	}

	return ""
}

func OperationNameToContext(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, operationNameContextKey, value)
}

func ClientTraceFromContext(ctx context.Context) bool {
	v := ctx.Value(clientTraceContextKey)
	if v != nil {
		return v.(bool)
	}

	return true
}

func ClientTraceToContext(ctx context.Context, value bool) context.Context {
	return context.WithValue(ctx, clientTraceContextKey, value)
}

func TraceRequest(tr opentracing.Tracer, req *http.Request, options ...nethttp.ClientOption) (*http.Request, *nethttp.Tracer) {
	opts := []nethttp.ClientOption{
		nethttp.ClientTrace(ClientTraceFromContext(req.Context())),
	}

	componentName := ComponentNameFromContext(req.Context())
	if componentName != "" {
		opts = append(opts, nethttp.ComponentName(componentName))
	}

	operationName := OperationNameFromContext(req.Context())
	if operationName != "" {
		opts = append(opts, nethttp.OperationName(operationName))
	}

	opts = append(opts, options...)

	return nethttp.TraceRequest(tr, req, opts...)
}
