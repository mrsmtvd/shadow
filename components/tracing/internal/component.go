package internal

import (
	"net"
	"strings"
	"sync"

	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/logging"
	"github.com/kihamo/shadow/components/metrics"
	"github.com/kihamo/shadow/components/tracing"
	"github.com/kihamo/shadow/components/tracing/internal/tracer"
	"github.com/opentracing/opentracing-go"
	jconfig "github.com/uber/jaeger-client-go/config"
)

type Component struct {
	application shadow.Application
	config      config.Component

	metricsOnce    sync.Once
	metricsFactory *factory
}

func (c *Component) Name() string {
	return tracing.ComponentName
}

func (c *Component) Version() string {
	return tracing.ComponentVersion
}

func (c *Component) Dependencies() []shadow.Dependency {
	return []shadow.Dependency{
		{
			Name:     config.ComponentName,
			Required: true,
		},
		{
			Name: logging.ComponentName,
		},
		{
			Name: metrics.ComponentName,
		},
	}
}

func (c *Component) Init(a shadow.Application) error {
	c.application = a
	c.config = a.GetComponent(config.ComponentName).(config.Component)

	return nil
}

func (c *Component) Run() error {
	return c.initTracer()
}

func (c *Component) initTracer() error {
	if !c.config.Bool(tracing.ConfigEnabled) {
		tracer.DefaultTracer.SetTracerNoop()

		return nil
	}

	tags := []opentracing.Tag{{
		Key:   tracing.TagAppVersion,
		Value: c.application.Version(),
	}, {
		Key:   tracing.TagAppBuild,
		Value: c.application.Build(),
	}}

	tagsFromConfig := c.config.String(tracing.ConfigTags)
	if len(tagsFromConfig) > 0 {
		var parts []string

		for _, tag := range strings.Split(tagsFromConfig, ",") {
			parts = strings.Split(tag, "=")

			if len(parts) > 1 {
				tags = append(tags, opentracing.Tag{
					Key:   strings.TrimSpace(parts[0]),
					Value: strings.TrimSpace(parts[1]),
				})
			}
		}
	}

	cfg := jconfig.Configuration{
		Disabled:    false,
		ServiceName: c.config.String(tracing.ConfigServiceName),
		Sampler: &jconfig.SamplerConfig{
			Type:                    c.config.String(tracing.ConfigSamplerType),
			Param:                   c.config.Float64(tracing.ConfigSamplerParam),
			SamplingServerURL:       c.config.String(tracing.ConfigSamplerServerURL),
			MaxOperations:           c.config.Int(tracing.ConfigSamplerMaxOperations),
			SamplingRefreshInterval: c.config.Duration(tracing.ConfigSamplerSamplingRefreshInterval),
		},
		Reporter: &jconfig.ReporterConfig{
			QueueSize:           c.config.Int(tracing.ConfigReporterQueueSize),
			BufferFlushInterval: c.config.Duration(tracing.ConfigReporterBufferFlushInterval),
			LogSpans:            c.config.Bool(tracing.ConfigReporterLogSpans),

			// Local collector
			LocalAgentHostPort: net.JoinHostPort(c.config.String(tracing.ConfigCollectorLocalHost), c.config.String(tracing.ConfigCollectorLocalPort)),

			// Remote collector
			CollectorEndpoint: c.config.String(tracing.ConfigCollectorRemoteEndpoint),
			User:              c.config.String(tracing.ConfigCollectorRemoteUser),
			Password:          c.config.String(tracing.ConfigCollectorRemotePassword),
		},
		Tags: tags,
	}

	options := make([]jconfig.Option, 0, 0)

	if c.application.HasComponent(logging.ComponentName) {
		log := logging.DefaultLogger().Named(c.Name())
		options = append(options, jconfig.Logger(NewLogger(log)))
	}

	if c.application.HasComponent(metrics.ComponentName) {
		options = append(options, jconfig.Metrics(c.newMetricsFactory()))
		cfg.RPCMetrics = c.config.Bool(tracing.ConfigMetricsRPCEnabled)
	}

	t, _, err := cfg.NewTracer(options...)
	if err != nil {
		return err
	}

	tracer.DefaultTracer.SetTracer(t)

	return nil
}

func (c *Component) Tracer() opentracing.Tracer {
	return tracer.DefaultTracer
}
