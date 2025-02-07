package interceptor

import (
	"fmt"
	"runtime/debug"

	"github.com/mrsmtvd/shadow/components/grpc/stats"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func NewRecoverUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = status.Errorf(codes.Internal, "%s", r)
				doRecover(ctx, info.FullMethod, req, err)
			}
		}()

		return handler(ctx, req)
	}
}

func NewRecoverStreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = status.Errorf(codes.Internal, "%s", r)
				doRecover(ss.Context(), info.FullMethod, nil, err)
			}
		}()

		return handler(srv, ss)
	}
}

func doRecover(ctx context.Context, fullMethod string, req interface{}, err error) {
	fields := map[string]interface{}{
		"error": err.Error(),
		"trace": string(debug.Stack()),
	}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if clientName := md.Get(stats.MetaDataClientNameKey); len(clientName) > 0 {
			fields["client"] = clientName[0]
		}
	}

	service, method := stats.SplitFullMethod(fullMethod)
	fields["service"] = service
	fields["method"] = method

	if req != nil {
		if s, ok := req.(fmt.Stringer); ok {
			fields["request"] = s
		}
	}

	grpclog.Error("Recovery panic", fields)
}
