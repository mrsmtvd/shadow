package internal

import (
	"net"

	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/logger"
	"github.com/kihamo/shadow/components/tracing"
	"github.com/opentracing/opentracing-go"
	jconfig "github.com/uber/jaeger-client-go/config"
)

type Component struct {
	application shadow.Application
	config      config.Component

	tracer *Tracer
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
			Name: logger.ComponentName,
		},
	}
}

func (c *Component) Init(a shadow.Application) error {
	c.application = a
	c.config = a.GetComponent(config.ComponentName).(config.Component)
	c.tracer = NewTracer()

	opentracing.SetGlobalTracer(c.tracer)

	return nil
}

func (c *Component) Run() error {
	return c.initTracer()
}

func (c *Component) initTracer() error {
	if !c.config.Bool(tracing.ConfigEnabled) {
		c.tracer.SetTracerNoop()

		return nil
	}

	cfg := jconfig.Configuration{
		Disabled:    false,
		ServiceName: c.application.Name(),
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
	}

	options := make([]jconfig.Option, 0, 0)

	if c.application.HasComponent(logger.ComponentName) {
		log := logger.NewOrNop(c.Name(), c.application)
		options = append(options, jconfig.Logger(NewLogger(log)))
	}

	tracer, _, err := cfg.NewTracer(options...)
	if err != nil {
		return err
	}

	c.tracer.SetTracer(tracer)

	return nil
}

func (c *Component) Tracer() opentracing.Tracer {
	return c.tracer
}
