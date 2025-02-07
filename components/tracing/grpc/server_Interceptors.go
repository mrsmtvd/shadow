package grpc

import (
	"context"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/mrsmtvd/shadow/components/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UnaryServerInterceptor(tracer opentracing.Tracer, opts ...otgrpc.Option) grpc.UnaryServerInterceptor {
	interceptor := otgrpc.OpenTracingServerInterceptor(tracer, opts...)

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		wrapped := func(ctx context.Context, req interface{}) (interface{}, error) {
			span := opentracing.SpanFromContext(ctx)
			if sc, ok := span.Context().(jaeger.SpanContext); ok {
				trailer := metadata.Pairs(tracing.ResponseTraceIDHeader, sc.TraceID().String())
				_ = grpc.SetTrailer(ctx, trailer)
			}

			return handler(ctx, req)
		}

		return interceptor(ctx, req, info, wrapped)
	}
}

func StreamServerInterceptor(tracer opentracing.Tracer, opts ...otgrpc.Option) grpc.StreamServerInterceptor {
	return otgrpc.OpenTracingStreamServerInterceptor(tracer, opts...)
}
