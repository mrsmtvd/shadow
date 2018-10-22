package tracer

import (
	"sync"

	"github.com/opentracing/opentracing-go"
)

type Wrapper struct {
	mutex  sync.RWMutex
	tracer opentracing.Tracer
}

func NewWrapper() *Wrapper {
	t := &Wrapper{}
	t.SetTracerNoop()

	return t
}

func (t *Wrapper) SetTracerNoop() {
	t.SetTracer(opentracing.NoopTracer{})
}

func (t *Wrapper) SetTracer(tr opentracing.Tracer) {
	t.mutex.Lock()
	t.tracer = tr
	t.mutex.Unlock()
}

func (t *Wrapper) StartSpan(operationName string, opts ...opentracing.StartSpanOption) opentracing.Span {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.tracer.StartSpan(operationName, opts...)
}

func (t *Wrapper) Inject(sm opentracing.SpanContext, format interface{}, carrier interface{}) error {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.tracer.Inject(sm, format, carrier)
}

func (t *Wrapper) Extract(format interface{}, carrier interface{}) (opentracing.SpanContext, error) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.tracer.Extract(format, carrier)
}
