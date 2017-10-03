package grpc

import (
	"time"

	"golang.org/x/net/context"
)

func (s *Server) Version(ctx context.Context, in *VersionRequest) (*VersionResponse, error) {
	return &VersionResponse{
		Name:      s.Application.GetName(),
		Version:   s.Application.GetVersion(),
		Build:     s.Application.GetBuild(),
		BuildDate: s.Application.GetBuildDate().Format(time.RFC3339),
		Uptime:    uint64(s.Application.GetUptime()),
	}, nil
}
