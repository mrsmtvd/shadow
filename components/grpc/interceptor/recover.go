package interceptor

import (
	"fmt"
	"runtime/debug"

	"github.com/kihamo/shadow/components/logger"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewRecoverUnaryServerInterceptor(l logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = status.Errorf(codes.Internal, "%s", r)
				doRecover(l, info.FullMethod, req, err)
			}
		}()

		return handler(ctx, req)
	}
}

func NewRecoverStreamServerInterceptor(l logger.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = status.Errorf(codes.Internal, "%s", r)
				doRecover(l, info.FullMethod, nil, err)
			}
		}()

		return handler(srv, ss)
	}
}

func NewRecoverUnaryClientInterceptor(l logger.Logger) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return nil
	}
}

func NewRecoverStreamClientInterceptor(l logger.Logger) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		return nil, nil
	}
}

func doRecover(l logger.Logger, method string, req interface{}, err error) {
	fields := map[string]interface{}{
		"method": method,
		"error":  err.Error(),
		"trace":  string(debug.Stack()),
	}

	if req != nil {
		if s, ok := req.(fmt.Stringer); ok {
			fields["request"] = s
		}
	}

	l.Error("Recovery panic", fields)
}
