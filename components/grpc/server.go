package grpc

import (
	g "google.golang.org/grpc"
	"google.golang.org/grpc/stats"
)

type HasGrpcServer interface {
	RegisterGrpcServer(s *g.Server)
}

type HasUnaryServerInterceptors interface {
	GrpcUnaryServerInterceptors() []g.UnaryServerInterceptor
}

type HasStreamServerInterceptors interface {
	GrpcStreamServerInterceptors() []g.StreamServerInterceptor
}

type HasStatsHandlers interface {
	GrpcStatsHandlers() []stats.Handler
}
