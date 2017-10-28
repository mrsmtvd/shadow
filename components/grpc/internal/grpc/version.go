package grpc

import (
	"github.com/golang/protobuf/ptypes"
	proto "github.com/kihamo/shadow/components/grpc/grpc"
	"golang.org/x/net/context"
)

func (s *Server) Version(ctx context.Context, in *proto.VersionRequest) (*proto.VersionResponse, error) {
	response := &proto.VersionResponse{
		Name:    s.Application.GetName(),
		Version: s.Application.GetVersion(),
		Build:   s.Application.GetBuild(),
		Uptime:  ptypes.DurationProto(s.Application.GetUptime()),
	}

	if s.Application.GetBuildDate() != nil {
		buildDatetime, err := ptypes.TimestampProto(*s.Application.GetBuildDate())
		if err != nil {
			return response, err
		}

		response.BuildDatetime = buildDatetime
	}

	if s.Application.GetStartDate() != nil {
		startDatetime, err := ptypes.TimestampProto(*s.Application.GetStartDate())
		if err != nil {
			return response, err
		}

		response.StartDatetime = startDatetime
	}

	return response, nil
}
