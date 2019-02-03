package logging

import (
	"context"
)

type contextKey struct{}

var (
	loggerContextKey = contextKey{}
)

func ContextWithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey, logger)
}

func LoggerFromContext(ctx context.Context) Logger {
	v := ctx.Value(loggerContextKey)
	if v != nil {
		return v.(Logger)
	}

	return DefaultLogger()
}

func Log(ctx context.Context) Logger {
	return LoggerFromContext(ctx)
}
