package server

import (
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/grpc/interceptor"
	"github.com/kihamo/shadow/components/grpc/stats"
	"google.golang.org/grpc"
	s "google.golang.org/grpc/stats"
)

func NewServer(opts ...grpc.ServerOption) *grpc.Server {
	return grpc.NewServer(opts...)
}

func NewServerWithDefaultOptions(config config.Component, opts ...grpc.ServerOption) *grpc.Server {
	opts = append([]grpc.ServerOption{
		WithDefaultStatsHandlerChain(config),
		WithDefaultUnaryChain(),
		WithDefaultStreamChain(),
	}, opts...)

	return grpc.NewServer(opts...)
}

func WithDefaultStatsHandlerChain(config config.Component, handlers ...s.Handler) grpc.ServerOption {
	if config != nil {
		handlers = append(handlers, stats.NewContextHandler(config))
	}

	handlers = append(handlers, []s.Handler{
		stats.NewLoggerHandler(),
		stats.NewMetricHandler(),
	}...)

	return WithStatsHandlerChain(handlers...)
}

func WithStatsHandlerChain(handlers ...s.Handler) grpc.ServerOption {
	return stats.WithStatsHandlerServerChain(handlers...)
}

func WithDefaultUnaryChain(interceptors ...grpc.UnaryServerInterceptor) grpc.ServerOption {
	interceptors = append(interceptors, interceptor.NewRecoverUnaryServerInterceptor())

	return WithUnaryChain(interceptors...)
}

func WithUnaryChain(interceptors ...grpc.UnaryServerInterceptor) grpc.ServerOption {
	return grpc_middleware.WithUnaryServerChain(interceptors...)
}

func WithDefaultStreamChain(interceptors ...grpc.StreamServerInterceptor) grpc.ServerOption {
	interceptors = append(interceptors, interceptor.NewRecoverStreamServerInterceptor())

	return WithStreamChain(interceptors...)
}

func WithStreamChain(interceptors ...grpc.StreamServerInterceptor) grpc.ServerOption {
	return grpc_middleware.WithStreamServerChain(interceptors...)
}
