package internal

import (
	proto "github.com/kihamo/shadow/components/grpc/grpc"
	handlers "github.com/kihamo/shadow/components/grpc/internal/grpc"
	"google.golang.org/grpc"
)

func (c *Component) RegisterGrpcServer(s *grpc.Server) {
	proto.RegisterGrpcServer(s, &handlers.Server{
		Application: c.application,
	})
}
