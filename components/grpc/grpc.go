package grpc

import (
	"google.golang.org/grpc"
)

type grpcServer struct {
	component *Component
}

func (c *Component) RegisterGrpcServer(s *grpc.Server) {
	RegisterGrpcServer(s, &grpcServer{
		component: c,
	})
}
