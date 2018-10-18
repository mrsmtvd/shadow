package tracing

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

func SpanError(span opentracing.Span, err error) {
	span.LogFields(log.Error(err))
	ext.Error.Set(span, true)
}

func StartSpanFromContext(ctx context.Context, componentName, operationName string, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
	tag := opentracing.Tag{Key: string(ext.Component), Value: componentName}
	opts = append([]opentracing.StartSpanOption{tag}, opts...)

	return opentracing.StartSpanFromContext(ctx, operationName, opts...)
}
