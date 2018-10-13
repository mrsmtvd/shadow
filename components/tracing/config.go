package tracing

const (
	ConfigEnabled                     = ComponentName + ".enabled"
	ConfigCollectorLocalHost          = ComponentName + ".collector.local.host"
	ConfigCollectorLocalPort          = ComponentName + ".collector.local.port"
	ConfigCollectorRemoteUser         = ComponentName + ".collector.remote.user"
	ConfigCollectorRemotePassword     = ComponentName + ".collector.remote.password"
	ConfigCollectorRemoteEndpoint     = ComponentName + ".collector.remote.endpoint"
	ConfigReporterQueueSize           = ComponentName + ".reporter.queue-size"
	ConfigReporterBufferFlushInterval = ComponentName + ".reporter.buffer-flush-interval"
	ConfigReporterLogSpans            = ComponentName + ".reporter.log-spans"
)
