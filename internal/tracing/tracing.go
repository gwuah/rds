package tracing

import (
	"context"
	"fmt"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

var tp *sdktrace.TracerProvider

const (
	tracerName       = "github.com/superfly/flyctl"
	HeaderFlyTraceId = "fly-trace-id"
	HeaderFlySpanId  = "fly-span-id"
)

func InitTraceProvider(ctx context.Context, appName string) (*sdktrace.TracerProvider, error) {
	if tp != nil {
		return tp, nil
	}

	var exporter sdktrace.SpanExporter
	switch {
	case os.Getenv("LOG_LEVEL") == "trace":
		stdoutExp, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			return nil, err
		}
		exporter = stdoutExp

	default:

		opts := []otlptracehttp.Option{
			otlptracehttp.WithEndpoint(getCollectorUrl() + ":4318"),
			otlptracehttp.WithInsecure(),
			otlptracehttp.WithHeaders(map[string]string{
				"authorization": getToken(ctx),
			}),
		}
		httpExporter, err := otlptracehttp.New(ctx, opts...)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to telemetry collector")
		}

		exporter = httpExporter
	}

	resourceAttrs := []attribute.KeyValue{
		semconv.ServiceNameKey.String("rds"),
		// attribute.String("build.info.version", buildinfo.Version().String()),
		// attribute.String("build.info.os", buildinfo.OS()),
		// attribute.String("build.info.arch", buildinfo.Arch()),
		// attribute.String("build.info.commit", buildinfo.Commit()),
	}

	if appName != "" {
		resourceAttrs = append(resourceAttrs, attribute.String("app.name", appName))
	}

	resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		resourceAttrs...,
	)

	tp = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource),
	)

	otel.SetTracerProvider(tp)

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	// otel.SetLogger(otelLogger(ctx))
	// otel.SetErrorHandler(errorHandler(ctx))

	return tp, nil
}

func getToken(ctx context.Context) string {
	token := ""
	if token == "" {
		token = os.Getenv("FLY_API_TOKEN")
	}
	return token
}
func getCollectorUrl() string {
	url := os.Getenv("FLY_TRACE_COLLECTOR_URL")
	if url != "" {
		return url
	}

	return "fly-otel-collector-dev.fly.dev"
}
