package tracing

import (
	"github.com/kihamo/shadow/components/tracing/internal/tracer"
	"github.com/opentracing/opentracing-go"
)

func DefaultTracer() opentracing.Tracer {
	return tracer.DefaultTracer
}
