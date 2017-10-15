package grpc

import (
	"github.com/gogo/protobuf/types"
	proto "github.com/kihamo/shadow/components/grpc/grpc"
	"golang.org/x/net/context"
)

func (s *Server) Version(ctx context.Context, in *proto.VersionRequest) (*proto.VersionResponse, error) {
	response := &proto.VersionResponse{
		Name:    s.Application.GetName(),
		Version: s.Application.GetVersion(),
		Build:   s.Application.GetBuild(),
		Uptime:  types.DurationProto(s.Application.GetUptime()),
	}

	if s.Application.GetBuildDate() != nil {
		buildDatetime, err := types.TimestampProto(*s.Application.GetBuildDate())
		if err != nil {
			return response, err
		}

		response.BuildDatetime = buildDatetime
	}

	return response, nil
}
