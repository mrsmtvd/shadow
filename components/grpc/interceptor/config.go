package interceptor

import (
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/kihamo/shadow/components/config"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var ConfigContextKey = &ContextKey{"config"}

func ConfigFromContext(c context.Context) *config.Component {
	return c.Value(ConfigContextKey).(*config.Component)
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

func NewConfigUnaryClientInterceptor(c config.Component) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return nil
	}
}

func NewConfigStreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		return nil, nil
	}
}
