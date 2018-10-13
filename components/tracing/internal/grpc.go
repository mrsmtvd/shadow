package internal

import (
	"context"

	"github.com/kihamo/shadow/components/tracing/grpc"
	g "google.golang.org/grpc"
)

func (c *Component) GrpcUnaryServerInterceptors() []g.UnaryServerInterceptor {
	return []g.UnaryServerInterceptor{
		func(ctx context.Context, req interface{}, info *g.UnaryServerInfo, handler g.UnaryHandler) (resp interface{}, err error) {
			return grpc.UnaryServerInterceptor(c.Tracer())(ctx, req, info, handler)
		},
	}
}

func (c *Component) GrpcStreamServerInterceptors() []g.StreamServerInterceptor {
	return []g.StreamServerInterceptor{
		func(srv interface{}, ss g.ServerStream, info *g.StreamServerInfo, handler g.StreamHandler) error {
			return grpc.StreamServerInterceptor(c.Tracer())(srv, ss, info, handler)
		},
	}
}
