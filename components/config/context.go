package config

import (
	"context"
)

type contextKey string

var (
	configContextKey = contextKey("config")
)

func ContextWithConfig(ctx context.Context, cfg Component) context.Context {
	return context.WithValue(ctx, configContextKey, cfg)
}

func FromContext(ctx context.Context) Component {
	v := ctx.Value(configContextKey)
	if v != nil {
		return v.(Component)
	}

	return nil
}
