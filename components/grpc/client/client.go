package client

import (
	"context"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/mrsmtvd/shadow/components/grpc/stats"
	"google.golang.org/grpc"
	s "google.golang.org/grpc/stats"
)

func Dial(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	return grpc.Dial(target, opts...)
}

func DefaultDial(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	return DefaultDialContext(context.Background(), target, opts...)
}

func DefaultDialWithCustomOptions(target string, unaryInterceptors []grpc.UnaryClientInterceptor, streamInterceptors []grpc.StreamClientInterceptor, statsHandlers []s.Handler, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	return DefaultDialContextWithCustomOptions(context.Background(), target, unaryInterceptors, streamInterceptors, statsHandlers, opts...)
}

func DialContext(ctx context.Context, target string, opts ...grpc.DialOption) (conn *grpc.ClientConn, err error) {
	return grpc.DialContext(ctx, target, opts...)
}

func DefaultDialContext(ctx context.Context, target string, opts ...grpc.DialOption) (conn *grpc.ClientConn, err error) {
	opts = append([]grpc.DialOption{
		WithDefaultUnaryChain(),
		WithDefaultStreamChain(),
		WithDefaultStatsHandlerChain(),
	}, opts...)

	return DefaultDialContextWithCustomOptions(ctx, target, nil, nil, nil, opts...)
}

func DefaultDialContextWithCustomOptions(ctx context.Context, target string, unaryInterceptors []grpc.UnaryClientInterceptor, streamInterceptors []grpc.StreamClientInterceptor, statsHandlers []s.Handler, opts ...grpc.DialOption) (conn *grpc.ClientConn, err error) {
	opts = append([]grpc.DialOption{
		WithDefaultUnaryChain(unaryInterceptors...),
		WithDefaultStreamChain(streamInterceptors...),
		WithDefaultStatsHandlerChain(statsHandlers...),
	}, opts...)

	return grpc.DialContext(ctx, target, opts...)
}

func WithDefaultStatsHandlerChain(handlers ...s.Handler) grpc.DialOption {
	return WithStatsHandlerChain(handlers...)
}

func WithStatsHandlerChain(handlers ...s.Handler) grpc.DialOption {
	return stats.WithStatsHandlerClientChain(handlers...)
}

func WithDefaultUnaryChain(interceptors ...grpc.UnaryClientInterceptor) grpc.DialOption {
	return WithUnaryChain(interceptors...)
}

func WithUnaryChain(interceptors ...grpc.UnaryClientInterceptor) grpc.DialOption {
	return grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(interceptors...))
}

func WithDefaultStreamChain(interceptors ...grpc.StreamClientInterceptor) grpc.DialOption {
	return WithStreamChain(interceptors...)
}

func WithStreamChain(interceptors ...grpc.StreamClientInterceptor) grpc.DialOption {
	return grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(interceptors...))
}
