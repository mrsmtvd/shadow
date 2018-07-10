package client

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/kihamo/shadow/components/grpc/stats"
	"github.com/kihamo/shadow/components/logger"
	"google.golang.org/grpc"
	s "google.golang.org/grpc/stats"
)

func Dial(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	return DialContext(context.Background(), target, opts...)
}

func DialContext(ctx context.Context, target string, opts ...grpc.DialOption) (conn *grpc.ClientConn, err error) {
	opts = append([]grpc.DialOption{WithDefaultStatsHandlerChain(nil)}, opts...)

	return grpc.DialContext(ctx, target, opts...)
}

func WithDefaultStatsHandlerChain(logger logger.Logger, handlers ...s.Handler) grpc.DialOption {
	if logger != nil {
		handlers = append(handlers, stats.NewLoggerHandler(logger))
	}

	handlers = append(handlers, stats.NewMetricHandler())

	return WithStatsHandlerChain(handlers...)
}

func WithStatsHandlerChain(handlers ...s.Handler) grpc.DialOption {
	return stats.WithStatsHandlerClientChain(handlers...)
}

func WithUnaryChain(interceptors ...grpc.UnaryClientInterceptor) grpc.DialOption {
	return grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(interceptors...))
}

func WithStreamChain(interceptors ...grpc.StreamClientInterceptor) grpc.DialOption {
	return grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(interceptors...))
}
