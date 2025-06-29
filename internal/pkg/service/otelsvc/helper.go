package otelsvc

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// StartSpan starts a new span and returns the new context and the span.
// It accepts a variadic list of trace.SpanStartOption to configure the span.
// If no trace.SpanKind is provided via the options, it defaults to trace.SpanKindServer.
// Any user-provided SpanKind option will override the default.
func StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	// Prepend the default span kind. User-provided options will be appended
	// and will override the default if they also set the span kind.
	finalOpts := append([]trace.SpanStartOption{trace.WithSpanKind(trace.SpanKindServer)}, opts...)

	ctx, span := otel.GetTracerProvider().Tracer(name).Start(
		ctx,
		name,
		finalOpts...,
	)

	return ctx, span
}

// StartSpanWithAttributes is a convenience wrapper around StartSpan that accepts attributes as a map[string]any.
// It safely converts the map values to OpenTelemetry attributes.
func StartSpanWithAttributes(ctx context.Context, name string, attrs map[string]any, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	var attributes []attribute.KeyValue
	for k, v := range attrs {
		switch val := v.(type) {
		case string:
			attributes = append(attributes, attribute.String(k, val))
		case int:
			attributes = append(attributes, attribute.Int(k, val))
		case int64:
			attributes = append(attributes, attribute.Int64(k, val))
		case bool:
			attributes = append(attributes, attribute.Bool(k, val))
		case float64:
			attributes = append(attributes, attribute.Float64(k, val))
		case []string:
			attributes = append(attributes, attribute.StringSlice(k, val))
		case []int:
			attributes = append(attributes, attribute.IntSlice(k, val))
		case []int64:
			attributes = append(attributes, attribute.Int64Slice(k, val))
		case []bool:
			attributes = append(attributes, attribute.BoolSlice(k, val))
		case []float64:
			attributes = append(attributes, attribute.Float64Slice(k, val))
		case uuid.UUID:
			attributes = append(attributes, attribute.String(k, val.String()))
		case []uuid.UUID:
			attributes = append(attributes, attribute.StringSlice(k, lo.Map(val, func(item uuid.UUID, _ int) string {
				return item.String()
			})))
		default:
			// For any other type, try to marshal it to a JSON string.
			if jsonBytes, err := json.Marshal(val); err == nil {
				attributes = append(attributes, attribute.String(k, string(jsonBytes)))
			}
			// If marshaling fails, the attribute is skipped.
		}
	}

	opts = append(opts, trace.WithAttributes(attributes...))

	return StartSpan(ctx, name, opts...)
}
