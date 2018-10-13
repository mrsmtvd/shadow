package internal

import (
	"context"

	"github.com/kihamo/shadow/components/tracing/grpc"
	g "google.golang.org/grpc"
)

func (c *Component) GrpcUnaryServerInterceptors() []g.UnaryServerInterceptor {
	return []g.UnaryServerInterceptor{
		func(ctx context.Context, req interface{}, info *g.UnaryServerInfo, handler g.UnaryHandler) (resp interface{}, err error) {
			tracer := c.Tracer()
			if tracer != nil {
				return grpc.UnaryServerInterceptor(tracer)(ctx, req, info, handler)
			}

			return handler(ctx, req)
		},
	}
}

func (c *Component) GrpcStreamServerInterceptors() []g.StreamServerInterceptor {
	return []g.StreamServerInterceptor{
		func(srv interface{}, ss g.ServerStream, info *g.StreamServerInfo, handler g.StreamHandler) error {
			tracer := c.Tracer()
			if tracer != nil {
				return grpc.StreamServerInterceptor(tracer)(srv, ss, info, handler)
			}

			return handler(srv, ss)
		},
	}
}
