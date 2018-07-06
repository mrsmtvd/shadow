package interceptor

import (
	"fmt"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging"
	"github.com/kihamo/shadow/components/logger"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

var LoggerContextKey = &ContextKey{"logger"}

func LoggerFromContext(c context.Context) logger.Logger {
	return c.Value(LoggerContextKey).(logger.Logger)
}

func NewLoggerUnaryServerInterceptor(l logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		resp, err = handler(context.WithValue(ctx, LoggerContextKey, l), req)

		doLogger(l, info.FullMethod, req, resp, err)

		return resp, err
	}
}

func NewLoggerStreamServerInterceptor(l logger.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		newCtx := context.WithValue(ss.Context(), LoggerContextKey, l)
		wrapped := grpc_middleware.WrapServerStream(ss)
		wrapped.WrappedContext = newCtx

		err := handler(srv, wrapped)

		doLogger(l, info.FullMethod, nil, nil, err)

		return err
	}
}

func NewLoggersUnaryClientInterceptor(l logger.Logger) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return nil
	}
}

func NewLoggerStreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		return nil, nil
	}
}

func doLogger(l logger.Logger, method string, req interface{}, resp interface{}, err error) {
	code := grpc_logging.DefaultErrorToCode(err)

	fields := map[string]interface{}{
		"method": method,
		"code":   code.String(),
	}

	if req != nil {
		if s, ok := req.(fmt.Stringer); ok {
			fields["request"] = s
		}
	}

	if resp != nil {
		if s, ok := resp.(fmt.Stringer); ok {
			fields["response"] = s
		}
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
