package telemetry

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"slices"

	"go.opentelemetry.io/contrib/bridges/otelzap"
	prombridge "go.opentelemetry.io/contrib/bridges/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var serviceName string

func StartSpan(ctx context.Context, name string) (context.Context, trace.Span) { //nolint:ireturn
	return otel.Tracer(serviceName).Start(ctx, name) //nolint:spancheck
}

// Signal is one OTLP signal. Each is pushed to wherever the standard
// OpenTelemetry environment points it — OTEL_EXPORTER_OTLP_<SIGNAL>_ENDPOINT
// (a full URL), or OTEL_EXPORTER_OTLP_ENDPOINT (a base URL the exporter
// appends /v1/<signal> to). A signal with neither set is not exported, so a
// process with no collector configured runs fine and quietly.
type Signal string

const (
	Traces  Signal = "TRACES"
	Metrics Signal = "METRICS"
	Logs    Signal = "LOGS"
)

// Enabled reports whether the signal has an endpoint configured.
func (s Signal) Enabled() bool {
	return os.Getenv("OTEL_EXPORTER_OTLP_"+string(s)+"_ENDPOINT") != "" ||
		os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") != ""
}

// Start starts the OTLP push pipeline for every configured signal and returns
// the shutdown that flushes and stops them. Everything is pushed — nothing
// is served for scraping — so it works the same in a long-lived process and
// in a function that is frozen between requests (pending data is flushed on
// the next thaw and by the shutdown):
//
//   - traces: the global tracer provider.
//   - metrics: the default Prometheus registry (request and pool collectors)
//     bridged into OTLP, plus anything recorded on the global meter.
//   - logs: telemetry.Log() is teed into OTLP; stdout logging stays.
func Start(ctx context.Context, appName string) (func(context.Context) error, error) {
	serviceName = appName
	otel.SetTextMapPropagator(propagation.TraceContext{})

	res, err := resource.New(ctx, resource.WithAttributes(semconv.ServiceName(appName), semconv.ServiceVersion("dev")))
	if err != nil {
		return nil, fmt.Errorf("telemetry resource: %w", err)
	}

	var shutdowns []func(context.Context) error
	shutdown := func(ctx context.Context) error {
		var errs []error
		for _, fn := range slices.Backward(shutdowns) {
			errs = append(errs, fn(ctx))
		}
		return errors.Join(errs...)
	}
	fail := func(err error) (func(context.Context) error, error) {
		_ = shutdown(ctx)
		return nil, err
	}

	if Traces.Enabled() {
		exporter, err := otlptracehttp.New(ctx)
		if err != nil {
			return fail(fmt.Errorf("trace exporter: %w", err))
		}
		tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter), sdktrace.WithResource(res))
		otel.SetTracerProvider(tp)
		shutdowns = append(shutdowns, tp.Shutdown)
	}

	if Metrics.Enabled() {
		exporter, err := otlpmetrichttp.New(ctx)
		if err != nil {
			return fail(fmt.Errorf("metric exporter: %w", err))
		}
		reader := sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithProducer(prombridge.NewMetricProducer()))
		mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader), sdkmetric.WithResource(res))
		otel.SetMeterProvider(mp)
		shutdowns = append(shutdowns, mp.Shutdown)
	}

	if Logs.Enabled() {
		exporter, err := otlploghttp.New(ctx)
		if err != nil {
			return fail(fmt.Errorf("log exporter: %w", err))
		}
		lp := sdklog.NewLoggerProvider(sdklog.WithProcessor(sdklog.NewBatchProcessor(exporter)), sdklog.WithResource(res))
		teeLogs(otelzap.NewCore(appName, otelzap.WithLoggerProvider(lp)))
		shutdowns = append(shutdowns, lp.Shutdown)
	}

	log.Printf("%s telemetry: traces=%t metrics=%t logs=%t", appName, Traces.Enabled(), Metrics.Enabled(), Logs.Enabled())
	return shutdown, nil
}

// teeLogs makes telemetry.Log() write to core as well as stdout.
func teeLogs(core zapcore.Core) {
	Init()
	logger = logger.WithOptions(zap.WrapCore(func(stdout zapcore.Core) zapcore.Core {
		return zapcore.NewTee(stdout, core)
	}))
}

// Shutdown flushes and stops the pipeline Start started, logging failures.
func Shutdown(ctx context.Context, shutdown func(context.Context) error, appName string) {
	if shutdown == nil {
		return
	}
	err := shutdown(ctx)
	if err != nil {
		Log().Error("error shutting down telemetry for "+appName, zap.Error(err))
	}
}
