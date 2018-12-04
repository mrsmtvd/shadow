package tracer

import (
	"github.com/opentracing/opentracing-go"
)

var DefaultTracer = NewWrapper()

func init() {
	opentracing.SetGlobalTracer(DefaultTracer)
}
