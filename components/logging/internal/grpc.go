package internal

import (
	cmp "github.com/mrsmtvd/shadow/components/grpc"
	"github.com/mrsmtvd/shadow/components/logging/grpc"
	"google.golang.org/grpc/stats"
)

func (c *Component) GrpcStatsHandlers() []stats.Handler {
	return []stats.Handler{
		grpc.NewStatsHandler(c.Logger().Named(cmp.ComponentName)),
	}
}
