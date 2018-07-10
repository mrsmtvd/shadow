package server

import (
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/grpc/interceptor"
	"github.com/kihamo/shadow/components/grpc/stats"
	"github.com/kihamo/shadow/components/logger"
	"google.golang.org/grpc"
	s "google.golang.org/grpc/stats"
)

func NewServer(opts ...grpc.ServerOption) *grpc.Server {
	return grpc.NewServer(opts...)
}

func NewServerWithDefaultServerOptions(config config.Component, logger logger.Logger, opts ...grpc.ServerOption) *grpc.Server {
	opts = append([]grpc.ServerOption{
		WithDefaultStatsHandlerChain(config, logger),
		WithDefaultUnaryChain(logger),
		WithDefaultStreamChain(logger),
	}, opts...)

	return grpc.NewServer(opts...)
}

func WithDefaultStatsHandlerChain(config config.Component, logger logger.Logger, handlers ...s.Handler) grpc.ServerOption {
	if config != nil {
		handlers = append(handlers, stats.NewContextHandler(config))
	}

	if logger != nil {
		handlers = append(handlers, stats.NewLoggerHandler(logger))
	}

	handlers = append(handlers, stats.NewMetricHandler())

	return WithStatsHandlerChain(handlers...)
}

func WithStatsHandlerChain(handlers ...s.Handler) grpc.ServerOption {
	return stats.WithStatsHandlerServerChain(handlers...)
}

func WithDefaultUnaryChain(logger logger.Logger, interceptors ...grpc.UnaryServerInterceptor) grpc.ServerOption {
	if logger != nil {
		interceptors = append(interceptors, interceptor.NewRecoverUnaryServerInterceptor(logger))
	}

	return WithUnaryChain(interceptors...)
}

func WithUnaryChain(interceptors ...grpc.UnaryServerInterceptor) grpc.ServerOption {
	return grpc_middleware.WithUnaryServerChain(interceptors...)
}

func WithDefaultStreamChain(logger logger.Logger, interceptors ...grpc.StreamServerInterceptor) grpc.ServerOption {
	if logger != nil {
		interceptors = append(interceptors, interceptor.NewRecoverStreamServerInterceptor(logger))
	}

	return WithStreamChain(interceptors...)
}

func WithStreamChain(interceptors ...grpc.StreamServerInterceptor) grpc.ServerOption {
	return grpc_middleware.WithStreamServerChain(interceptors...)
}
