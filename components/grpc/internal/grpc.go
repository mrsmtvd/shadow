package internal

import (
	handlers "github.com/mrsmtvd/shadow/components/grpc/internal/grpc"
	"github.com/mrsmtvd/shadow/components/grpc/protobuf"
	"google.golang.org/grpc"
)

func (c *Component) RegisterGrpcServer(s *grpc.Server) {
	protobuf.RegisterGrpcServer(s, &handlers.Server{
		Application: c.application,
	})
}
