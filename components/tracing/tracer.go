package tracing

import (
	"github.com/mrsmtvd/shadow/components/tracing/internal/tracer"
	"github.com/opentracing/opentracing-go"
)

func DefaultTracer() opentracing.Tracer {
	return tracer.DefaultTracer
}
