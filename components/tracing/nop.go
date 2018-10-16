package tracing

import (
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/tracing/internal"
	"github.com/opentracing/opentracing-go"
)

func NewOrNop(application shadow.Application) opentracing.Tracer {
	if cmp := application.GetComponent(ComponentName); cmp != nil {
		return cmp.(Component).Tracer()
	}

	return internal.NewTracer()
}
