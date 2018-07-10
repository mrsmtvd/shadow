package stats

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/stats"
)

type LoggerHandler struct {
	Handler
}

func NewLoggerHandler() *LoggerHandler {
	return &LoggerHandler{}
}

func (h *LoggerHandler) HandleConn(ctx context.Context, stat stats.ConnStats) {
	ctxValue := h.ConnectValueFromContext(ctx)
	fields := map[string]interface{}{
		"remote-address": ctxValue.RemoteAddr,
		"local-address":  ctxValue.LocalAddr,
	}

	switch stat.(type) {
	case *stats.ConnBegin:
		grpclog.Info("GRPC connect is open", fields)

	case *stats.ConnEnd:
		grpclog.Info("GRPC connect is closed", fields)
	}
}

func (h *LoggerHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	ctxValue := h.RPCValueFromContext(ctx)
	fields := map[string]interface{}{
		"service": ctxValue.Service,
		"method":  ctxValue.Method,
		"type":    ctxValue.Type,
		"client":  ctxValue.ClientName,
	}

	switch s := stat.(type) {
	case *stats.InPayload:
		if s.IsClient() {
			fields["response"] = s.Payload
		} else {
			fields["request"] = s.Payload
		}

		grpclog.Info("Call gRPC method", fields)

	case *stats.OutPayload:
		if s.IsClient() {
			fields["request"] = s.Payload
		} else {
			fields["response"] = s.Payload
		}

		grpclog.Info("Call gRPC method", fields)

	case *stats.End:
		code := grpc_logging.DefaultErrorToCode(s.Error)
		fields["code"] = code

		if s.Error != nil {
			fields["error"] = s.Error.Error()
		}

		switch code {
		case codes.OK:
			grpclog.Info("Called gRPC method", fields)
		case codes.OutOfRange, codes.Internal, codes.Unavailable, codes.DataLoss:
			grpclog.Error("Called gRPC method", fields)
		default:
			grpclog.Info("Called gRPC method", fields)
		}
	}
}
