package tracing

const (
	ConfigEnabled                        = ComponentName + ".enabled"
	ConfigServiceName                    = ComponentName + ".service-name"
	ConfigTags                           = ComponentName + ".tags"
	ConfigCollectorLocalHost             = ComponentName + ".collector.local.host"
	ConfigCollectorLocalPort             = ComponentName + ".collector.local.port"
	ConfigCollectorRemoteUser            = ComponentName + ".collector.remote.user"
	ConfigCollectorRemotePassword        = ComponentName + ".collector.remote.password"
	ConfigCollectorRemoteEndpoint        = ComponentName + ".collector.remote.endpoint"
	ConfigReporterQueueSize              = ComponentName + ".reporter.queue-size"
	ConfigReporterBufferFlushInterval    = ComponentName + ".reporter.buffer-flush-interval"
	ConfigReporterLogSpans               = ComponentName + ".reporter.log-spans"
	ConfigSamplerType                    = ComponentName + ".sampler.type"
	ConfigSamplerParam                   = ComponentName + ".sampler.param"
	ConfigSamplerServerURL               = ComponentName + ".sampler.server-url"
	ConfigSamplerMaxOperations           = ComponentName + ".sampler.max-operations"
	ConfigSamplerSamplingRefreshInterval = ComponentName + ".sampler.sampling-refresh-interval"
)
