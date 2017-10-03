package grpc

import (
	"google.golang.org/grpc"
)

type HasGrpcServer interface {
	RegisterGrpcServer(s *grpc.Server)
}
