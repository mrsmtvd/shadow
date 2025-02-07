package internal

import (
	"github.com/mrsmtvd/shadow/components/tracing/grpc"
	g "google.golang.org/grpc"
)

func (c *Component) GrpcUnaryServerInterceptors() []g.UnaryServerInterceptor {
	return []g.UnaryServerInterceptor{
		grpc.UnaryServerInterceptor(c.Tracer()),
	}
}

func (c *Component) GrpcStreamServerInterceptors() []g.StreamServerInterceptor {
	return []g.StreamServerInterceptor{
		grpc.StreamServerInterceptor(c.Tracer()),
	}
}
