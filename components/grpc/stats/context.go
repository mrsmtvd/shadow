package stats

import (
	"context"

	"github.com/kihamo/shadow/components/config"
	"google.golang.org/grpc/stats"
)

var (
	configContextKey = &contextKey{"config"}
)

func ConfigFromContext(c context.Context) config.Component {
	return c.Value(configContextKey).(config.Component)
}

type ContextHandler struct {
	Handler

	config config.Component
}

func NewContextHandler(c config.Component) *ContextHandler {
	return &ContextHandler{
		config: c,
	}
}

func (h *ContextHandler) TagConn(ctx context.Context, info *stats.ConnTagInfo) context.Context {
	return h.Handler.TagConn(context.WithValue(ctx, configContextKey, h.config), info)
}
