package grpc

import (
	"time"

	"golang.org/x/net/context"
)

func (s *grpcServer) Version(ctx context.Context, in *VersionRequest) (*VersionResponse, error) {
	return &VersionResponse{
		Name:      s.component.application.GetName(),
		Version:   s.component.application.GetVersion(),
		Build:     s.component.application.GetBuild(),
		BuildDate: s.component.application.GetBuildDate().Format(time.RFC3339),
		Uptime:    uint64(s.component.application.GetUptime()),
	}, nil
}
