package httpserver

import (
	"context"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"lib/telemetry"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// processStart is as early as Go code sees the process: package variables
// are set before main runs. The boot span starts here, so it covers the
// runtime's start and the telemetry setup that comes before any span.
var processStart = time.Now()

// placementEnv names the variables the Cloudflare container runtime sets
// about where an instance runs. Elsewhere they are unset.
var placementEnv = []struct{ key, env string }{
	{"cloudflare.location", "CLOUDFLARE_LOCATION"},
	{"cloudflare.region", "CLOUDFLARE_REGION"},
	{"cloudflare.country", "CLOUDFLARE_COUNTRY_A2"},
	{"cloudflare.durable_object_id", "CLOUDFLARE_DURABLE_OBJECT_ID"},
}

// placement is where this process runs, for the boot span.
func placement() []attribute.KeyValue {
	var attrs []attribute.KeyValue
	for _, p := range placementEnv {
		if v := os.Getenv(p.env); v != "" {
			attrs = append(attrs, attribute.String(p.key, v))
		}
	}
	return attrs
}

// startBoot opens the boot span: one trace per process start, from
// processStart until the port is bound. Its children are the startup before
// it (boot.init), each module, and Build; the database connect and the
// constraint check sit under the postgres module.
func startBoot(ctx context.Context, appName string) (context.Context, trace.Span) { //nolint:ireturn
	ctx, span := telemetry.StartSpan(ctx, "boot",
		trace.WithTimestamp(processStart),
		trace.WithAttributes(append(placement(), attribute.String("boot.app", appName))...))
	_, initSpan := telemetry.StartSpan(ctx, "boot.init", trace.WithTimestamp(processStart))
	initSpan.End()
	return ctx, span
}

// failBoot ends the boot span as failed.
func failBoot(span trace.Span, err error) {
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
	span.End()
}

// egressTraceURL answers with the Cloudflare location that received the
// request ("colo=GRU", "loc=BR"). From a container that is where its
// outbound traffic leaves Cloudflare, which decides how far every database
// round trip travels.
const egressTraceURL = "https://cloudflare.com/cdn-cgi/trace"

const egressProbeTimeout = 5 * time.Second

// probeEgress records where this process's outbound traffic leaves, as a
// boot.egress span under the boot trace. It runs after the port is bound,
// so it never delays a boot, and a failure only marks the span.
func probeEgress(ctx context.Context) {
	ctx, span := telemetry.StartSpan(ctx, "boot.egress")
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, egressProbeTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, egressTraceURL, nil)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return
	}
	for line := range strings.SplitSeq(string(body), "\n") {
		key, value, _ := strings.Cut(line, "=")
		switch key {
		case "colo":
			span.SetAttributes(attribute.String("egress.colo", value))
		case "loc":
			span.SetAttributes(attribute.String("egress.country", value))
		}
	}
}
