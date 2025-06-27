package otelsvc

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func StartSpan(ctx context.Context, name string, attributes ...map[string]string) (context.Context, trace.Span) {
	ctx, span := otel.GetTracerProvider().Tracer(name).Start(
		ctx,
		name,
		trace.WithSpanKind(trace.SpanKindServer),
	)
	for _, attr := range attributes {
		for k, v := range attr {
			span.SetAttributes(attribute.Key(k).String(v))
		}
	}

	return ctx, span
}
