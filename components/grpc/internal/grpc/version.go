package grpc

import (
	"fmt"
	"math/rand"

	"github.com/golang/protobuf/ptypes"
	proto "github.com/kihamo/shadow/components/grpc/grpc"
	"golang.org/x/net/context"
	g "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
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

	if true {
		header := metadata.Pairs("header-key", "val")
		g.SendHeader(ctx, header)

		r := rand.Intn(2)
		fmt.Println(r)

		if r == 1 {
			return nil, status.Error(codes.NotFound, "111")
		}
	}

	return response, nil
}
