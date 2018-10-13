package grpc

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging"
	st "github.com/kihamo/shadow/components/grpc/stats"
	"github.com/kihamo/shadow/components/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/stats"
)

type StatsHandler struct {
	st.Handler

	log logger.Logger
}

func NewStatsHandler(log logger.Logger) *StatsHandler {
	return &StatsHandler{
		log: log,
	}
}

func (h *StatsHandler) HandleConn(ctx context.Context, stat stats.ConnStats) {
	ctxValue := h.ConnectValueFromContext(ctx)
	fields := map[string]interface{}{
		"remote-address": ctxValue.RemoteAddr,
		"local-address":  ctxValue.LocalAddr,
	}

	switch stat.(type) {
	case *stats.ConnBegin:
		h.log.Info("GRPC connect is open", fields)

	case *stats.ConnEnd:
		h.log.Info("GRPC connect is closed", fields)
	}
}

func (h *StatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
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

		h.log.Info("Call gRPC method", fields)

	case *stats.OutPayload:
		if s.IsClient() {
			fields["request"] = s.Payload
		} else {
			fields["response"] = s.Payload
		}

		h.log.Info("Call gRPC method", fields)

	case *stats.End:
		code := grpc_logging.DefaultErrorToCode(s.Error)
		fields["code"] = code

		if s.Error != nil {
			fields["error"] = s.Error.Error()
		}

		switch code {
		case codes.OK:
			h.log.Info("Called gRPC method", fields)
		case codes.OutOfRange, codes.Internal, codes.Unavailable, codes.DataLoss:
			h.log.Error("Called gRPC method", fields)
		default:
			h.log.Info("Called gRPC method", fields)
		}
	}
}
