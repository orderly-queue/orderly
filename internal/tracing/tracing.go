package tracing

import (
	"context"

	"github.com/henrywhitaker3/go-template/internal/app"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	otrace "go.opentelemetry.io/otel/trace"
)

var (
	Tracer = otel.Tracer("api")
)

func InitTracer(app *app.App) (*trace.TracerProvider, error) {
	var exporter trace.SpanExporter
	var err error
	if app.Config.Telemetry.Tracing.Endpoint == "stdout" {
		exporter, err = stdout.New(stdout.WithPrettyPrint())
	} else {
		exporter, err = otlptracehttp.New(context.Background(), otlptracehttp.WithEndpointURL(app.Config.Telemetry.Tracing.Endpoint))
	}
	if err != nil {
		return nil, err
	}

	res, err := resource.New(
		context.Background(),
		resource.WithContainer(),
		resource.WithOS(),
		resource.WithHost(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithFromEnv(),
		resource.WithAttributes(
			semconv.ServiceName("api"),
			semconv.ServiceVersion(app.Version),
			attribute.String("environment", app.Config.Environment),
		),
	)
	if err != nil {
		return nil, err
	}

	tp := trace.NewTracerProvider(
		trace.WithSampler(trace.TraceIDRatioBased(app.Config.Telemetry.Tracing.SampleRate)),
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, nil
}

func NewSpan(ctx context.Context, name string, opts ...otrace.SpanStartOption) (context.Context, otrace.Span) {
	return Tracer.Start(ctx, name, opts...)
}

func AddString(ctx context.Context, key, value string) {
	if span := otrace.SpanFromContext(ctx); span != nil {
		span.SetAttributes(attribute.String(key, value))
	}
}

func TraceID(ctx context.Context) string {
	span := otrace.SpanFromContext(ctx)
	if span.SpanContext().HasTraceID() {
		return span.SpanContext().TraceID().String()
	}
	return ""
}
