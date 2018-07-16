package stats

import (
	"context"
	"net"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/stats"
)

type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "grpc context value " + k.name
}

var (
	connectContextKey = &contextKey{"connect"}
	rpcContextKey     = &contextKey{"rpc"}
)

type ConnectContextValue struct {
	RemoteAddr net.Addr
	LocalAddr  net.Addr
}

type RPCContextValue struct {
	Service    string
	Method     string
	Type       string
	ClientName string
}

func SplitFullMethod(name string) (string, string) {
	if i := strings.LastIndex(name, "/"); i >= 0 {
		return name[1:i], name[i+1:]
	}
	return "unknown", "unknown"
}

func CodeAsString(code codes.Code) string {
	switch code {
	case codes.OK:
		return "OK"
	case codes.Canceled:
		return "CANCELLED"
	case codes.InvalidArgument:
		return "INVALID_ARGUMENT"
	case codes.DeadlineExceeded:
		return "DEADLINE_EXCEEDED"
	case codes.NotFound:
		return "NOT_FOUND"
	case codes.AlreadyExists:
		return "ALREADY_EXISTS"
	case codes.PermissionDenied:
		return "PERMISSION_DENIED"
	case codes.ResourceExhausted:
		return "RESOURCE_EXHAUSTED"
	case codes.FailedPrecondition:
		return "FAILED_PRECONDITION"
	case codes.Aborted:
		return "ABORTED"
	case codes.OutOfRange:
		return "OUT_OF_RANGE"
	case codes.Unimplemented:
		return "UNIMPLEMENTED"
	case codes.Internal:
		return "INTERNAL"
	case codes.Unavailable:
		return "UNAVAILABLE"
	case codes.DataLoss:
		return "DATA_LOSS"
	case codes.Unauthenticated:
		return "UNAUTHENTICATED"
	default:
		return "UNKNOWN"
	}
}

type Handler struct {
	stats.Handler
}

func (h Handler) ConnectValueFromContext(ctx context.Context) *ConnectContextValue {
	v := ctx.Value(connectContextKey)

	if v == nil {
		return nil
	}

	if value, ok := v.(*ConnectContextValue); ok {
		return value
	}

	return nil
}

func (h Handler) RPCValueFromContext(ctx context.Context) *RPCContextValue {
	v := ctx.Value(rpcContextKey)

	if v == nil {
		return nil
	}

	if value, ok := v.(*RPCContextValue); ok {
		return value
	}

	return nil
}

func (h *Handler) TagConn(ctx context.Context, info *stats.ConnTagInfo) context.Context {
	if exist := h.ConnectValueFromContext(ctx); exist != nil {
		return ctx
	}

	ctxValue := &ConnectContextValue{
		RemoteAddr: info.RemoteAddr,
		LocalAddr:  info.LocalAddr,
	}

	return context.WithValue(ctx, connectContextKey, ctxValue)
}

func (h *Handler) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	if exist := h.RPCValueFromContext(ctx); exist != nil {
		return ctx
	}

	ctxValue := &RPCContextValue{
		ClientName: DefaultClientName,
	}

	ctxValue.Service, ctxValue.Method = SplitFullMethod(info.FullMethodName)

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if clientName := md.Get(MetaDataClientNameKey); len(clientName) > 0 {
			ctxValue.ClientName = clientName[0]
		}
	}

	return context.WithValue(ctx, rpcContextKey, ctxValue)
}

func (h *Handler) HandleConn(ctx context.Context, stat stats.ConnStats) {

}

func (h *Handler) HandleRPC(ctx context.Context, stat stats.RPCStats) {

}
