package telemetry

import (
	"context"

	"lib/errx"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.uber.org/zap"
)

func InitTracer(ctx context.Context, appName string) func(context.Context) error {
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpointURL("http://victoria-traces:10428/insert/opentelemetry/v1/traces"),
	)
	if err != nil {
		errx.Exit(err, "failed to init tracer for "+appName)
	}

	res, err := resource.New(ctx, resource.WithAttributes(semconv.ServiceName(appName), semconv.ServiceVersion("dev")))
	if err != nil {
		errx.Exit(err, "failed to init tracer for "+appName)
	}

	tp := trace.NewTracerProvider(trace.WithBatcher(exporter), trace.WithResource(res))
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tp.Shutdown
}

func ShutdownTracer(ctx context.Context, shutdown func(context.Context) error, appName string) {
	if err := shutdown(ctx); err != nil {
		Log().Error("error shutting down tracer for "+appName, zap.Error(err))
	}
}
