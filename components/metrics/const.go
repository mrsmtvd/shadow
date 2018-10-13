package metrics

const (
	ComponentName    = "metrics"
	ComponentVersion = "3.0.0"

	ProtocolHTTP = "http"

	StatusOK      = "ok"
	StatusTimeout = "timeout"
	StatusError   = "error"

	// For all responses
	MetricNameResponseTimeSeconds        = ComponentName + "_response_time_seconds"
	MetricNameResponseSizeBytes          = ComponentName + "_response_size_bytes"
	MetricNameResponseMarshalTimeSeconds = ComponentName + "_response_marshal_time_seconds"

	// For all requests
	MetricNameRequestSizeBytes                = ComponentName + "_request_size_bytes"
	MetricNameRequestsTotal                   = ComponentName + "_requests_total"
	MetricNameRequestReadTimeSeconds          = ComponentName + "_requests_read_time_seconds"
	MetricNameRequestReadUnmarshalTimeSeconds = ComponentName + "_requests_read_unmarshal_time_seconds"
	MetricNameRequestUnmarshalTimeSeconds     = ComponentName + "_requests_unmarshal_time_seconds"

	// For all external requests
	MetricNameExternalResponseTimeSeconds = ComponentName + "_external_response_time_seconds"

	// For GRPC
	MetricNameHandledTotal  = ComponentName + "_grpc_handled_total"
	MetricNameReceivedTotal = ComponentName + "_grpc_received_total"
	MetricNameSentTotal     = ComponentName + "_grpc_sent_total"
	MetricNameStartedTotal  = ComponentName + "_grpc_started_total"
)
