package grpc

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/kihamo/shadow/components/grpc/protobuf"
	"golang.org/x/net/context"
)

func (s *Server) Version(ctx context.Context, in *protobuf.VersionRequest) (*protobuf.VersionResponse, error) {
	response := &protobuf.VersionResponse{
		Name:    s.Application.Name(),
		Version: s.Application.Version(),
		Build:   s.Application.Build(),
		Uptime:  ptypes.DurationProto(s.Application.Uptime()),
	}

	if s.Application.BuildDate() != nil {
		buildDatetime, err := ptypes.TimestampProto(*s.Application.BuildDate())
		if err != nil {
			return response, err
		}

		response.BuildDatetime = buildDatetime
	}

	if s.Application.StartDate() != nil {
		startDatetime, err := ptypes.TimestampProto(*s.Application.StartDate())
		if err != nil {
			return response, err
		}

		response.StartDatetime = startDatetime
	}

	return response, nil
}
