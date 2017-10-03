package grpc

import (
	"google.golang.org/grpc"
)

type HasGrpcService interface {
	RegisterGrpcServer(s *grpc.Server)
}
