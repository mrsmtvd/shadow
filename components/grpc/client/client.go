package client

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/kihamo/shadow/components/grpc/stats"
	"google.golang.org/grpc"
	s "google.golang.org/grpc/stats"
)

func DialWithDefaultOptions(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	return DialContextWithDefaultOptions(context.Background(), target, opts...)
}

func DialContextWithDefaultOptions(ctx context.Context, target string, opts ...grpc.DialOption) (conn *grpc.ClientConn, err error) {
	opts = append([]grpc.DialOption{WithDefaultStatsHandlerChain()}, opts...)

	return grpc.DialContext(ctx, target, opts...)
}

func WithDefaultStatsHandlerChain(handlers ...s.Handler) grpc.DialOption {
	handlers = append(handlers, []s.Handler{
		stats.NewLoggerHandler(),
		stats.NewMetricHandler(),
	}...)

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
