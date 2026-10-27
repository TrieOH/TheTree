package httpserver

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

// TestTraceMiddlewareNamesSpanAfterOperation — one span per operation,
// named by its operationID, visible to the handler through both the ctx
// and the request (AuthDispatch reads the request's context), with a
// returned error recorded on it.
func TestTraceMiddlewareNamesSpanAfterOperation(t *testing.T) {
	rec := tracetest.NewSpanRecorder()
	prev := otel.GetTracerProvider()
	otel.SetTracerProvider(sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(rec)))
	t.Cleanup(func() { otel.SetTracerProvider(prev) })

	boom := errors.New("boom")
	handler := TraceMiddleware()(func(ctx context.Context, _ http.ResponseWriter, r *http.Request, _ any) (any, error) {
		if !trace.SpanFromContext(ctx).SpanContext().IsValid() {
			t.Error("handler ctx carries no span")
		}
		if !trace.SpanFromContext(r.Context()).SpanContext().Equal(trace.SpanFromContext(ctx).SpanContext()) {
			t.Error("request context does not carry the operation span")
		}
		return nil, boom
	}, "CreateStep")

	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/forms/x/steps", nil)
	_, err := handler(req.Context(), httptest.NewRecorder(), req, nil)
	if !errors.Is(err, boom) {
		t.Fatalf("want the handler error back, got %v", err)
	}

	spans := rec.Ended()
	if len(spans) != 1 {
		t.Fatalf("want 1 span, got %d", len(spans))
	}
	if spans[0].Name() != "CreateStep" {
		t.Fatalf("want span named CreateStep, got %q", spans[0].Name())
	}
	if spans[0].Status().Code != codes.Error {
		t.Fatalf("want error status, got %v", spans[0].Status().Code)
	}
}
