package grpc

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging"
	st "github.com/kihamo/shadow/components/grpc/stats"
	"github.com/kihamo/shadow/components/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/stats"
)

type StatsHandler struct {
	st.Handler

	log logging.Logger
}

func NewStatsHandler(log logging.Logger) *StatsHandler {
	return &StatsHandler{
		log: log,
	}
}

func (h *StatsHandler) HandleConn(ctx context.Context, stat stats.ConnStats) {
	ctxValue := h.ConnectValueFromContext(ctx)
	fields := []interface{}{
		"remote-address", ctxValue.RemoteAddr.String(),
		"local-address", ctxValue.LocalAddr.String(),
	}

	switch stat.(type) {
	case *stats.ConnBegin:
		h.log.Info("GRPC connect is open", fields...)

	case *stats.ConnEnd:
		h.log.Info("GRPC connect is closed", fields...)
	}
}

func (h *StatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	ctxValue := h.RPCValueFromContext(ctx)
	fields := []interface{}{
		"service", ctxValue.Service,
		"method", ctxValue.Method,
		"type", ctxValue.Type,
		"client", ctxValue.ClientName,
	}

	switch s := stat.(type) {
	case *stats.InPayload:
		if s.IsClient() {
			fields = append(fields, "response", s.Payload)
		} else {
			fields = append(fields, "request", s.Payload)
		}

		h.log.Info("Call gRPC method", fields...)

	case *stats.OutPayload:
		if s.IsClient() {
			fields = append(fields, "request", s.Payload)
		} else {
			fields = append(fields, "response", s.Payload)
		}

		h.log.Info("Call gRPC method", fields...)

	case *stats.End:
		code := grpc_logging.DefaultErrorToCode(s.Error)
		fields = append(fields, "code", code)

		if s.Error != nil {
			fields = append(fields, "error", s.Error.Error())
		}

		switch code {
		case codes.OK:
			h.log.Info("Called gRPC method", fields...)
		case codes.OutOfRange, codes.Internal, codes.Unavailable, codes.DataLoss:
			h.log.Error("Called gRPC method", fields...)
		default:
			h.log.Info("Called gRPC method", fields...)
		}
	}
}
