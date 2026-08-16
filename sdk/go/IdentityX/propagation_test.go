package idx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// TestTraceContextPropagation verifies that calls made under an active span
// inject a traceparent header, so the receiving API continues the same
// distributed trace instead of starting a fresh root trace.
func TestTraceContextPropagation(t *testing.T) {
	// A real tracer provider so spans carry valid IDs — otelhttp only
	// injects traceparent when the span context is valid.
	prevTP := otel.GetTracerProvider()
	prevProp := otel.GetTextMapPropagator()
	otel.SetTracerProvider(sdktrace.NewTracerProvider())
	otel.SetTextMapPropagator(propagation.TraceContext{})
	t.Cleanup(func() {
		otel.SetTracerProvider(prevTP)
		otel.SetTextMapPropagator(prevProp)
	})

	var got trace.SpanContext
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = trace.SpanContextFromContext(
			propagation.TraceContext{}.Extract(context.Background(), propagation.HeaderCarrier(r.Header)),
		)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"keys": []}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(Config{BaseURL: srv.URL, APIKey: "test-api-key", ProjectID: uuid.New()})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	ctx, span := otel.Tracer("propagation-test").Start(context.Background(), "root")
	defer span.End()

	if _, err := client.Tokens.GetJWKS(ctx, true); err != nil {
		t.Fatalf("GetJWKS: %v", err)
	}

	want := span.SpanContext()
	if !got.IsValid() {
		t.Fatal("server received no valid trace context (traceparent header missing)")
	}
	if !got.IsRemote() {
		t.Fatal("received span context should be marked remote")
	}
	// The header carries the transport's client span, so only the trace ID
	// must match the caller's span — that is what joins the two services
	// into one distributed trace.
	if got.TraceID() != want.TraceID() {
		t.Fatalf("trace id = %s, want %s (trace not continued across services)", got.TraceID(), want.TraceID())
	}
	if got.SpanID() == want.SpanID() {
		t.Fatal("expected the transport to create its own client span")
	}
}
