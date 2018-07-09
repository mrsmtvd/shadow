package client

import (
	"context"

	"github.com/kihamo/shadow/components/grpc/stats"
	"github.com/kihamo/shadow/components/logger"
	"google.golang.org/grpc"
)

func Dial(target string, logger logger.Logger, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	return DialContext(context.Background(), target, logger, opts...)
}

func DialContext(ctx context.Context, target string, logger logger.Logger, opts ...grpc.DialOption) (conn *grpc.ClientConn, err error) {
	opts = append([]grpc.DialOption{
		stats.WithStatsHandlerClientChain(
			stats.NewLoggerHandler(logger),
			stats.NewMetricHandler()),
	}, opts...)

	return grpc.DialContext(ctx, target, opts...)
}
