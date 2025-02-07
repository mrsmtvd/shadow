package tracing

import (
	"github.com/mrsmtvd/shadow"
	"github.com/opentracing/opentracing-go"
)

type Component interface {
	shadow.Component

	Tracer() opentracing.Tracer
}
