package internal

import (
	"github.com/mrsmtvd/shadow/components/metrics/grpc"
	"google.golang.org/grpc/stats"
)

func (c *Component) GrpcStatsHandlers() []stats.Handler {
	return []stats.Handler{
		grpc.NewStatsHandler(),
	}
}
