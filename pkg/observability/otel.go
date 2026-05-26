package observability

import (
	"context"
	"fmt"
	"os"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.40.0"
)

// InitTracing configures OpenTelemetry tracing and returns a shutdown function.
func InitTracing(ctx context.Context) (func(context.Context) error, error) {
	exporterName := strings.ToLower(getEnv("OTEL_TRACES_EXPORTER", "stdout"))
	if exporterName == "none" {
		return func(context.Context) error { return nil }, nil
	}

	exporter, err := newTraceExporter(ctx, exporterName)
	if err != nil {
		return nil, err
	}

	serviceName := getEnv("OTEL_SERVICE_NAME", "mauna-backend")
	serviceVersion := getEnv("OTEL_SERVICE_VERSION", "2.0.0")

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(serviceVersion),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create telemetry resource: %w", err)
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return provider.Shutdown, nil
}

func newTraceExporter(ctx context.Context, exporterName string) (sdktrace.SpanExporter, error) {
	switch exporterName {
	case "otlp":
		return otlptracehttp.New(ctx)
	case "stdout":
		return stdouttrace.New(stdouttrace.WithPrettyPrint())
	default:
		return nil, fmt.Errorf("unsupported OTEL_TRACES_EXPORTER %q", exporterName)
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
