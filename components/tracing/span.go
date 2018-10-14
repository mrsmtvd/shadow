package tracing

import (
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

func SpanError(span opentracing.Span, err error) {
	span.LogFields(log.Error(err))
	ext.Error.Set(span, true)
}
