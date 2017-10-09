package grpc

import (
	"time"

	proto "github.com/kihamo/shadow/components/grpc/grpc"
	"golang.org/x/net/context"
)

func (s *Server) Version(ctx context.Context, in *proto.VersionRequest) (*proto.VersionResponse, error) {
	return &proto.VersionResponse{
		Name:      s.Application.GetName(),
		Version:   s.Application.GetVersion(),
		Build:     s.Application.GetBuild(),
		BuildDate: s.Application.GetBuildDate().Format(time.RFC3339),
		Uptime:    uint64(s.Application.GetUptime()),
	}, nil
}
