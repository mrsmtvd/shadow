package internal

import (
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
	}
}

func (c *Component) ConfigWatchers() []config.Watcher {
	return []config.Watcher{
		config.NewWatcher([]string{
			tracing.ConfigEnabled,
			tracing.ConfigCollectorLocalHost,
			tracing.ConfigCollectorLocalPort,
			tracing.ConfigCollectorRemoteUser,
			tracing.ConfigCollectorRemotePassword,
			tracing.ConfigCollectorRemoteEndpoint,
			tracing.ConfigReporterQueueSize,
			tracing.ConfigReporterBufferFlushInterval,
			tracing.ConfigReporterLogSpans,
		}, c.watchReInit),
	}
}

func (c *Component) watchReInit(_ string, newValue interface{}, _ interface{}) {
	c.initTracer()
}
