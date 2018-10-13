package internal

import (
	"sync"

	"github.com/opentracing/opentracing-go"
)

type Tracer struct {
	mutex  sync.Mutex
	tracer opentracing.Tracer
}

func NewTracer() *Tracer {
	t := &Tracer{}
	t.SetTracerNoop()

	return t
}

func (t *Tracer) SetTracerNoop() {
	t.SetTracer(opentracing.NoopTracer{})
}

func (t *Tracer) SetTracer(tracer opentracing.Tracer) {
	t.mutex.Lock()
	t.tracer = tracer
	t.mutex.Unlock()
}

func (t *Tracer) StartSpan(operationName string, opts ...opentracing.StartSpanOption) opentracing.Span {
	return t.tracer.StartSpan(operationName, opts...)
}

func (t *Tracer) Inject(sm opentracing.SpanContext, format interface{}, carrier interface{}) error {
	return t.tracer.Inject(sm, format, carrier)
}

func (t *Tracer) Extract(format interface{}, carrier interface{}) (opentracing.SpanContext, error) {
	return t.tracer.Extract(format, carrier)
}
