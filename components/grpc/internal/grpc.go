package internal

import (
	"github.com/kihamo/shadow/components/grpc/internal/grpc"
	g "google.golang.org/grpc"
)

func (c *Component) RegisterGrpcServer(s *g.Server) {
	grpc.RegisterGrpcServer(s, &grpc.Server{
		Application: c.application,
	})
}
