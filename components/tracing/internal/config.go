package internal

import (
	"strings"
	"time"

	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/tracing"
	"github.com/uber/jaeger-client-go"
)

func (c *Component) ConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariable(tracing.ConfigEnabled, config.ValueTypeBool).
			WithUsage("Enabled").
			WithGroup("Others").
			WithDefault(false).
			WithEditable(true),
		config.NewVariable(tracing.ConfigServiceName, config.ValueTypeString).
			WithUsage("Service name").
			WithGroup("Others").
			WithDefaultFunc(func() interface{} {
				if c.application == nil {
					return "shadow"
				}

				serviceName := c.application.Name()
				serviceName = strings.ToLower(serviceName)
				serviceName = strings.Join(strings.Fields(serviceName), "-")
				return serviceName
			}).
			WithEditable(true),
		config.NewVariable(tracing.ConfigCollectorLocalHost, config.ValueTypeString).
			WithUsage("Host").
			WithGroup("Local collector").
			WithDefault(jaeger.DefaultUDPSpanServerHost).
			WithEditable(true),
		config.NewVariable(tracing.ConfigCollectorLocalPort, config.ValueTypeInt).
			WithUsage("Port number").
			WithGroup("Local collector").
			WithDefault(jaeger.DefaultUDPSpanServerPort).
			WithEditable(true),
		config.NewVariable(tracing.ConfigCollectorRemoteUser, config.ValueTypeString).
			WithUsage("User login").
			WithGroup("Remote collector").
			WithEditable(true),
		config.NewVariable(tracing.ConfigCollectorRemotePassword, config.ValueTypeString).
			WithUsage("User password").
			WithGroup("Remote collector").
			WithEditable(true).
			WithView([]string{config.ViewPassword}),
		config.NewVariable(tracing.ConfigCollectorRemoteEndpoint, config.ValueTypeString).
			WithUsage("Endpoint").
			WithGroup("Remote collector").
			WithEditable(true),
		config.NewVariable(tracing.ConfigReporterQueueSize, config.ValueTypeInt).
			WithUsage("Size of internal queue where reported spans are stored before they are processed in the background").
			WithGroup("Reporter").
			WithDefault(100).
			WithEditable(true),
		config.NewVariable(tracing.ConfigReporterBufferFlushInterval, config.ValueTypeDuration).
			WithUsage("How often the buffer is force-flushed, even if it's not full").
			WithGroup("Reporter").
			WithDefault(time.Second).
			WithEditable(true),
		config.NewVariable(tracing.ConfigReporterLogSpans, config.ValueTypeBool).
			WithUsage("When enabled, enables LoggingReporter that runs in parallel with the main reporter and logs all submitted spans").
			WithGroup("Reporter").
			WithDefault(false).
			WithEditable(true),
		config.NewVariable(tracing.ConfigSamplerType, config.ValueTypeString).
			WithUsage("Type").
			WithGroup("Sampler").
			WithDefault(jaeger.SamplerTypeRemote).
			WithEditable(true).
			WithView([]string{config.ViewEnum}).
			WithViewOptions(map[string]interface{}{
				config.ViewOptionEnumOptions: [][]interface{}{
					{jaeger.SamplerTypeConst, "Constant"},
					{jaeger.SamplerTypeProbabilistic, "Probabilistic"},
					{jaeger.SamplerTypeRateLimiting, "Rate limiting"},
					// {jaeger.SamplerTypeLowerBound, "Lower bound"},
					{jaeger.SamplerTypeRemote, "Remote"},
				},
			}),
		config.NewVariable(tracing.ConfigSamplerParam, config.ValueTypeFloat64).
			WithUsage("Param is a value passed to the sampler").
			WithGroup("Sampler").
			WithDefault(0).
			WithEditable(true),
		config.NewVariable(tracing.ConfigSamplerServerURL, config.ValueTypeString).
			WithUsage("Address of jaeger-agent's HTTP sampling server (only for remote sampler)").
			WithGroup("Sampler").
			WithDefault("http://localhost:5778/sampling").
			WithEditable(true),
		config.NewVariable(tracing.ConfigSamplerMaxOperations, config.ValueTypeInt).
			WithUsage("Maximum number of operations that the sampler will keep track of. If an operation is not tracked, a default probabilistic sampler will be used rather than the per operation specific sampler").
			WithGroup("Sampler").
			WithDefault(2000).
			WithEditable(true),
		config.NewVariable(tracing.ConfigSamplerSamplingRefreshInterval, config.ValueTypeDuration).
			WithUsage("Controls how often the remotely controlled sampler will poll jaeger-agent for the appropriate sampling strategy").
			WithGroup("Sampler").
			WithDefault(time.Minute).
			WithEditable(true),
	}
}

func (c *Component) ConfigWatchers() []config.Watcher {
	return []config.Watcher{
		config.NewWatcher([]string{
			tracing.ConfigEnabled,
			tracing.ConfigServiceName,
			tracing.ConfigCollectorLocalHost,
			tracing.ConfigCollectorLocalPort,
			tracing.ConfigCollectorRemoteUser,
			tracing.ConfigCollectorRemotePassword,
			tracing.ConfigCollectorRemoteEndpoint,
			tracing.ConfigReporterQueueSize,
			tracing.ConfigReporterBufferFlushInterval,
			tracing.ConfigReporterLogSpans,
			tracing.ConfigSamplerType,
			tracing.ConfigSamplerParam,
			tracing.ConfigSamplerServerURL,
			tracing.ConfigSamplerMaxOperations,
			tracing.ConfigSamplerSamplingRefreshInterval,
		}, c.watchReInit),
	}
}

func (c *Component) watchReInit(_ string, newValue interface{}, _ interface{}) {
	c.initTracer()
}
