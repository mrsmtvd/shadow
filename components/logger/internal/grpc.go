package internal

import (
	"github.com/kihamo/shadow/components/logger/grpc"
	"google.golang.org/grpc/stats"
)

func (c *Component) GrpcStatsHandlers() []stats.Handler {
	return []stats.Handler{
		grpc.NewStatsHandler(c.Get("grpc")),
	}
}
