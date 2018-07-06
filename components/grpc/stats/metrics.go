package stats

import (
	"context"
	"fmt"
	"strings"

	"github.com/kihamo/shadow/components/grpc"
	"github.com/kihamo/snitch"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"
)

const (
	metaDataClientNameKey = "user-agent"

	// GRPC specific
	MetricNameHandledTotal  = grpc.ComponentName + "_handled_total"
	MetricNameReceivedTotal = grpc.ComponentName + "_received_total"
	MetricNameSentTotal     = grpc.ComponentName + "_sent_total"
	MetricNameStartedTotal  = grpc.ComponentName + "_started_total"

	// For all requests
	MetricNameResponseTimeSeconds = grpc.ComponentName + "_response_time_seconds"
	MetricNameRequestSizeBytes    = grpc.ComponentName + "_request_size_bytes"
	MetricNameResponseSizeBytes   = grpc.ComponentName + "_response_size_bytes"
	MetricNameRequestsTotal       = grpc.ComponentName + "_requests_total"
)

var (
	connectGRPCContextKey = &contextKey{"connect"}

	MetricHandledTotal  = snitch.NewCounter(MetricNameHandledTotal, "GRPC handled requests total")
	MetricReceivedTotal = snitch.NewCounter(MetricNameReceivedTotal, "GRPC received requests total")
	MetricSentTotal     = snitch.NewCounter(MetricNameSentTotal, "GRPC sent responses total")
	MetricStartedTotal  = snitch.NewCounter(MetricNameStartedTotal, "GRPC started requests total")

	MetricResponseTimeSeconds = snitch.NewTimer(MetricNameResponseTimeSeconds, "GRPC response time in total")
	MetricRequestSizeBytes    = snitch.NewHistogram(MetricNameRequestSizeBytes, "GRPC size of requests in bytes")
	MetricResponseSizeBytes   = snitch.NewHistogram(MetricNameResponseSizeBytes, "GRPC size of responses in bytes")
	MetricRequestsTotal       = snitch.NewCounter(MetricNameRequestsTotal, "GRPC requests total")
)

type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "grpc context value " + k.name
}

type ConnectGRPCContextValue struct {
	Service    string
	Method     string
	Type       string
	ClientName string
}

type MetricsHandler struct {
	stats.Handler
}

func NewMetricHandler() *MetricsHandler {
	return &MetricsHandler{}
}

func (h *MetricsHandler) TagConn(ctx context.Context, info *stats.ConnTagInfo) context.Context {
	//fmt.Println("TagConn::")

	return ctx
}

func (h *MetricsHandler) HandleConn(ctx context.Context, stat stats.ConnStats) {
	//fmt.Println("HandleConn")

	/*
		switch s := stat.(type) {
		case *stats.ConnBegin:
			fmt.Println("HandleConn::ConnBegin", s.IsClient())

		case *stats.ConnEnd:
			fmt.Println("HandleConn::ConnEnd")

		default:
			fmt.Println("HandleConn::default")
		}
	*/
}

func (h *MetricsHandler) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	//fmt.Println("TagRPC::")

	ctxValue := ConnectGRPCContextValue{
		Type: "", // TODO:
	}

	ctxValue.Service, ctxValue.Method = split(info.FullMethodName)

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if clientName := md.Get(metaDataClientNameKey); len(clientName) > 0 {
			ctxValue.ClientName = clientName[0]
		}
	}

	return context.WithValue(ctx, connectGRPCContextKey, ctxValue)
}

func (h *MetricsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	v := ctx.Value(connectGRPCContextKey)
	if v == nil {
		return
	}
	ctxValue := v.(ConnectGRPCContextValue)

	switch s := stat.(type) {
	case *stats.Begin:
		//fmt.Println("HandleRPC::Begin")

		MetricStartedTotal.With(
			"grpc_service", ctxValue.Service,
			"grpc_method", ctxValue.Method,
			"grpc_type", ctxValue.Type,
			"client_name", ctxValue.ClientName).Inc()

		MetricRequestsTotal.With(
			"handler", fmt.Sprintf("%s/%s", ctxValue.Service, ctxValue.Method),
			"protocol", grpc.ProtocolGRPC,
			"client_name", ctxValue.ClientName).Inc()

	case *stats.End:
		//fmt.Println("HandleRPC::End")

		if !s.IsClient() {
			code, _ := status.FromError(s.Error)
			MetricHandledTotal.With(
				"grpc_service", ctxValue.Service,
				"grpc_method", ctxValue.Method,
				"grpc_type", ctxValue.Type,
				"client_name", ctxValue.ClientName,
				"grpc_code", codeAsString(code.Code())).Inc()

			st := grpc.StatusOK
			if s.Error != nil {
				st = grpc.StatusError
			}
			MetricResponseTimeSeconds.With(
				"handler", fmt.Sprintf("%s/%s", ctxValue.Service, ctxValue.Method),
				"protocol", grpc.ProtocolGRPC,
				"client_name", ctxValue.ClientName,
				"status", st).UpdateSince(s.BeginTime)
		}

	case *stats.InPayload:
		//fmt.Println("HandleRPC::InPayload")

		if !s.IsClient() {
			MetricReceivedTotal.With(
				"grpc_service", ctxValue.Service,
				"grpc_method", ctxValue.Method,
				"grpc_type", ctxValue.Type,
				"client_name", ctxValue.ClientName).Inc()

			MetricRequestSizeBytes.With(
				"handler", fmt.Sprintf("%s/%s", ctxValue.Service, ctxValue.Method),
				"protocol", grpc.ProtocolGRPC,
				"client_name", ctxValue.ClientName,
			).Add(float64(s.WireLength))
		}

	case *stats.OutPayload:
		//fmt.Println("HandleRPC::OutPayload")

		if !s.IsClient() {
			MetricSentTotal.With(
				"grpc_service", ctxValue.Service,
				"grpc_method", ctxValue.Method,
				"grpc_type", ctxValue.Type,
				"client_name", ctxValue.ClientName).Inc()

			MetricResponseSizeBytes.With(
				"handler", fmt.Sprintf("%s/%s", ctxValue.Service, ctxValue.Method),
				"protocol", grpc.ProtocolGRPC,
				"client_name", ctxValue.ClientName,
				"status", grpc.StatusOK).Add(float64(s.WireLength))
		}

	case *stats.InTrailer:
		//fmt.Println("HandleRPC::InTrailer")

	case *stats.OutTrailer:
		//fmt.Println("HandleRPC::OutTrailer")

	case *stats.InHeader:
		//fmt.Println("HandleRPC::InHeader")

	case *stats.OutHeader:
		//fmt.Println("HandleRPC::OutHeader")

	default:
		//fmt.Println("HandleRPC::default")
	}
}

func codeAsString(code codes.Code) string {
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

func split(name string) (string, string) {
	if i := strings.LastIndex(name, "/"); i >= 0 {
		return name[1:i], name[i+1:]
	}
	return "unknown", "unknown"
}
