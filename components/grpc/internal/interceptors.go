package internal

import (
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/logger"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

var (
	ConfigContextKey = &ContextKey{"config"}
	LoggerContextKey = &ContextKey{"logger"}
)

type ContextKey struct {
	name string
}

func (k *ContextKey) String() string {
	return "grpc context value " + k.name
}
func ConfigFromContext(c context.Context) *config.Component {
	return c.Value(ConfigContextKey).(*config.Component)
}

func LoggerFromContext(c context.Context) logger.Logger {
	return c.Value(LoggerContextKey).(logger.Logger)
}

func NewConfigUnaryServerInterceptor(c config.Component) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		newCtx := context.WithValue(ctx, ConfigContextKey, c)

		return handler(newCtx, req)
	}
}

func NewConfigStreamServerInterceptor(c config.Component) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		newCtx := context.WithValue(ss.Context(), ConfigContextKey, c)
		wrapped := grpc_middleware.WrapServerStream(ss)
		wrapped.WrappedContext = newCtx

		return handler(srv, wrapped)
	}
}

func doLogger(l logger.Logger, method string, err error) {
	code := grpc_logging.DefaultErrorToCode(err)

	fields := map[string]interface{}{
		"method": method,
		"code":   code.String(),
	}

	if err != nil {
		fields["error"] = err.Error()

		l.Error("Called unary gRPC method", fields)
	} else {
		switch code {
		case codes.OK:
			l.Debug("Called unary gRPC method", fields)
		case codes.OutOfRange, codes.Internal, codes.Unavailable, codes.DataLoss:
			l.Error("Called unary gRPC method", fields)
		default:
			l.Warn("Called unary gRPC method", fields)
		}
	}
}

func NewLoggerUnaryServerInterceptor(l logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		resp, err = handler(context.WithValue(ctx, LoggerContextKey, l), req)

		doLogger(l, info.FullMethod, err)

		return resp, err
	}
}

func NewLoggerStreamServerInterceptor(l logger.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		newCtx := context.WithValue(ss.Context(), LoggerContextKey, l)
		wrapped := grpc_middleware.WrapServerStream(ss)
		wrapped.WrappedContext = newCtx

		err := handler(srv, wrapped)

		doLogger(l, info.FullMethod, err)

		return err
	}
}

func doMetrics(method string, startTime time.Time) {
	if metricExecuteTime != nil {
		metricExecuteTime.UpdateSince(startTime)
		metricExecuteTime.With("method", method).UpdateSince(startTime)
	}
}

func NewMetricsUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		startTime := time.Now()

		resp, err = handler(ctx, req)

		doMetrics(info.FullMethod, startTime)

		return resp, err
	}
}

func NewMetricsStreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		startTime := time.Now()

		err := handler(srv, ss)

		doMetrics(info.FullMethod, startTime)

		return err
	}
}
