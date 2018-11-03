package internal

import (
	cmp "github.com/kihamo/shadow/components/grpc"
	"github.com/kihamo/shadow/components/logging"
	"github.com/kihamo/shadow/components/logging/grpc"
	"google.golang.org/grpc/stats"
)

func (c *Component) GrpcStatsHandlers() []stats.Handler {
	return []stats.Handler{
		grpc.NewStatsHandler(logging.DefaultLogger().Named(cmp.ComponentName)),
	}
}
