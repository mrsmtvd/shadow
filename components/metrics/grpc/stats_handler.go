package grpc

import (
	"context"
	"fmt"

	"github.com/kihamo/shadow/components/grpc"
	st "github.com/kihamo/shadow/components/grpc/stats"
	"github.com/kihamo/shadow/components/metrics"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"
)

type StatsHandler struct {
	st.Handler
}

func NewStatsHandler() *StatsHandler {
	return &StatsHandler{}
}

func (h *StatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	ctxValue := h.RPCValueFromContext(ctx)

	switch s := stat.(type) {
	case *stats.Begin:
		if !s.IsClient() {
			metrics.MetricGRPCStartedTotal.With(
				"grpc_service", ctxValue.Service,
				"grpc_method", ctxValue.Method,
				"grpc_type", ctxValue.Type,
				"client_name", ctxValue.ClientName).Inc()

			metrics.MetricRequestsTotal.With(
				"handler", fmt.Sprintf("%s/%s", ctxValue.Service, ctxValue.Method),
				"protocol", grpc.ProtocolGRPC,
				"client_name", ctxValue.ClientName).Inc()
		}

	case *stats.End:
		responseTime := s.EndTime.Sub(s.BeginTime)
		sts := status.Convert(s.Error)

		code := metrics.StatusOK
		if sts.Code() == codes.DeadlineExceeded {
			code = metrics.StatusTimeout
		} else if s.Error != nil {
			code = metrics.StatusError
		}

		if !s.IsClient() {
			metrics.MetricGRPCHandledTotal.With(
				"grpc_service", ctxValue.Service,
				"grpc_method", ctxValue.Method,
				"grpc_type", ctxValue.Type,
				"client_name", ctxValue.ClientName,
				"grpc_code", st.CodeAsString(sts.Code())).Inc()

			metrics.MetricResponseTimeSeconds.With(
				"handler", fmt.Sprintf("%s/%s", ctxValue.Service, ctxValue.Method),
				"protocol", grpc.ProtocolGRPC,
				"client_name", ctxValue.ClientName,
				"status", code).Update(responseTime)
		} else {
			metrics.MetricExternalResponseTimeSeconds.With(
				"external_service", ctxValue.Service,
				"method", ctxValue.Method,
				"status", code).Update(responseTime)
		}

	case *stats.InPayload:
		if !s.IsClient() {
			metrics.MetricGRPCReceivedTotal.With(
				"grpc_service", ctxValue.Service,
				"grpc_method", ctxValue.Method,
				"grpc_type", ctxValue.Type,
				"client_name", ctxValue.ClientName).Inc()

			metrics.MetricRequestSizeBytes.With(
				"handler", fmt.Sprintf("%s/%s", ctxValue.Service, ctxValue.Method),
				"protocol", grpc.ProtocolGRPC,
				"client_name", ctxValue.ClientName,
			).Add(float64(s.WireLength))
		}

	case *stats.OutPayload:
		if !s.IsClient() {
			metrics.MetricGRPCSentTotal.With(
				"grpc_service", ctxValue.Service,
				"grpc_method", ctxValue.Method,
				"grpc_type", ctxValue.Type,
				"client_name", ctxValue.ClientName).Inc()

			metrics.MetricResponseSizeBytes.With(
				"handler", fmt.Sprintf("%s/%s", ctxValue.Service, ctxValue.Method),
				"protocol", grpc.ProtocolGRPC,
				"client_name", ctxValue.ClientName,
				"status", metrics.StatusOK).Add(float64(s.WireLength))
		}

	case *stats.InTrailer:

	case *stats.OutTrailer:

	case *stats.InHeader:

	case *stats.OutHeader:

	default:
	}
}
