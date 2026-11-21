package telemetry

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

func TestMaskIP(t *testing.T) {
	cases := map[string]string{
		"138.117.221.225":                         "138.117.221.0/24",
		"138.117.221.225:443":                     "138.117.221.0/24",
		"::ffff:10.1.2.3":                         "10.1.2.0/24",
		"2804:299c:110a:e400:cc96:8d63:9ee4:9ca8": "2804:299c:110a::/48",
		"[2804:1b3:7000:ce52::1]:8080":            "2804:1b3:7000::/48",
		"not-an-ip":                               "",
		"":                                        "",
	}
	for in, want := range cases {
		if got := MaskIP(in); got != want {
			t.Errorf("MaskIP(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestIPMaskerMasksClientAddressesOnSpans(t *testing.T) {
	rec := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(ipMasker{}), sdktrace.WithSpanProcessor(rec))
	_, span := tp.Tracer("test").Start(context.Background(), "req",
		trace.WithAttributes(
			attribute.String("client.address", "168.197.58.158"),
			attribute.String("network.peer.address", "2a06:98c0:3600::103"),
			attribute.String("http.route", "/events"),
		))
	span.End()

	got := map[attribute.Key]string{}
	for _, kv := range rec.Ended()[0].Attributes() {
		got[kv.Key] = kv.Value.AsString()
	}
	if got["client.address"] != "168.197.58.0/24" {
		t.Errorf("client.address = %q", got["client.address"])
	}
	if got["network.peer.address"] != "2a06:98c0:3600::/48" {
		t.Errorf("network.peer.address = %q", got["network.peer.address"])
	}
	if got["http.route"] != "/events" {
		t.Errorf("unrelated attribute changed: http.route = %q", got["http.route"])
	}
}
