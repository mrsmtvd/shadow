package stats

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging"
	"github.com/kihamo/shadow/components/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/stats"
)

type LoggerHandler struct {
	Handler

	logger logger.Logger
}

func NewLoggerHandler(l logger.Logger) *LoggerHandler {
	return &LoggerHandler{
		logger: l,
	}
}

func (h *LoggerHandler) HandleConn(ctx context.Context, stat stats.ConnStats) {
	ctxValue := h.ConnectValueFromContext(ctx)
	fields := map[string]interface{}{
		"remote-address": ctxValue.RemoteAddr,
		"local-address":  ctxValue.LocalAddr,
	}

	switch stat.(type) {
	case *stats.ConnBegin:
		h.logger.Debug("GRPC connect is open", fields)

	case *stats.ConnEnd:
		h.logger.Debug("GRPC connect is closed", fields)
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

		h.logger.Debug("Call gRPC method", fields)

	case *stats.OutPayload:
		if s.IsClient() {
			fields["request"] = s.Payload
		} else {
			fields["response"] = s.Payload
		}

		h.logger.Debug("Call gRPC method", fields)

	case *stats.End:
		code := grpc_logging.DefaultErrorToCode(s.Error)
		fields["code"] = code

		if s.Error != nil {
			fields["error"] = s.Error.Error()
		}

		switch code {
		case codes.OK:
			h.logger.Debug("Called gRPC method", fields)
		case codes.OutOfRange, codes.Internal, codes.Unavailable, codes.DataLoss:
			h.logger.Error("Called gRPC method", fields)
		default:
			h.logger.Warn("Called gRPC method", fields)
		}
	}
}
