package server

import (
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/kihamo/shadow/components/grpc/interceptor"
	"github.com/kihamo/shadow/components/grpc/stats"
	"google.golang.org/grpc"
	s "google.golang.org/grpc/stats"
)

func NewServer(opts ...grpc.ServerOption) *grpc.Server {
	return grpc.NewServer(opts...)
}

func NewDefaultServer(opts ...grpc.ServerOption) *grpc.Server {
	opts = append([]grpc.ServerOption{
		WithDefaultUnaryChain(),
		WithDefaultStreamChain(),
		WithDefaultStatsHandlerChain(),
	}, opts...)

	return NewDefaultServerWithCustomOptions(nil, nil, nil)
}

func NewDefaultServerWithCustomOptions(unaryInterceptors []grpc.UnaryServerInterceptor, streamInterceptors []grpc.StreamServerInterceptor, statsHandlers []s.Handler, opts ...grpc.ServerOption) *grpc.Server {
	opts = append([]grpc.ServerOption{
		WithDefaultUnaryChain(unaryInterceptors...),
		WithDefaultStreamChain(streamInterceptors...),
		WithDefaultStatsHandlerChain(statsHandlers...),
	}, opts...)

	return grpc.NewServer(opts...)
}

func WithDefaultStatsHandlerChain(handlers ...s.Handler) grpc.ServerOption {
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
