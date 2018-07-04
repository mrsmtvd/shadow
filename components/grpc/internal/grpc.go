package internal

import (
	handlers "github.com/kihamo/shadow/components/grpc/internal/grpc"
	"github.com/kihamo/shadow/components/grpc/protobuf"
	"google.golang.org/grpc"
)

func (c *Component) RegisterGrpcServer(s *grpc.Server) {
	protobuf.RegisterGrpcServer(s, &handlers.Server{
		Application: c.application,
	})
}
