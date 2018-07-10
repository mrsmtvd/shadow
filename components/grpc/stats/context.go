package stats

import (
	"context"

	"github.com/kihamo/shadow/components/config"
	"google.golang.org/grpc/stats"
)

var (
	configContextKey = &contextKey{"config"}
)

func ConfigFromContext(ctx context.Context) config.Component {
	v := ctx.Value(configContextKey)

	if v == nil {
		return nil
	}

	if value, ok := v.(config.Component); ok {
		return value
	}

	return nil
}

func ConfigToContext(ctx context.Context, c config.Component) context.Context {
	return context.WithValue(ctx, configContextKey, c)
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
	return h.Handler.TagConn(ConfigToContext(ctx, h.config), info)
}
