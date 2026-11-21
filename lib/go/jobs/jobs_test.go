package jobs_test

import (
	"context"
	"errors"
	"testing"

	"lib/jobs"

	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

type greet struct {
	Name string `json:"name"`
}

func (greet) Kind() string { return "greet" }

func TestDispatchRoundTrip(t *testing.T) {
	reg := jobs.NewRegistry()
	var got string
	jobs.Register(reg, func(_ context.Context, j greet) error {
		got = j.Name
		return nil
	})

	env, err := jobs.Encode(context.Background(), greet{Name: "ada"})
	if err != nil {
		t.Fatal(err)
	}
	err = reg.Dispatch(context.Background(), env)
	if err != nil {
		t.Fatal(err)
	}
	if got != "ada" {
		t.Fatalf("handler got %q, want ada", got)
	}
}

func TestDispatchUnknownKind(t *testing.T) {
	err := jobs.NewRegistry().Dispatch(context.Background(), jobs.Envelope{Kind: "nope"})
	if !errors.Is(err, jobs.ErrUnknownKind) {
		t.Fatalf("err = %v, want ErrUnknownKind", err)
	}
}

func TestRegisterTwicePanics(t *testing.T) {
	reg := jobs.NewRegistry()
	jobs.Register(reg, func(context.Context, greet) error { return nil })
	defer func() {
		if recover() == nil {
			t.Fatal("second Register did not panic")
		}
	}()
	jobs.Register(reg, func(context.Context, greet) error { return nil })
}

// A job joins the trace that enqueued it, and links to the span that
// delivered it (the /_jobs request on Cloudflare).
func TestDispatchJoinsTheEnqueuersTrace(t *testing.T) {
	rec := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(rec))
	prev := otel.GetTracerProvider()
	otel.SetTracerProvider(tp)
	t.Cleanup(func() { otel.SetTracerProvider(prev) })
	tracer := tp.Tracer("test")

	reqCtx, request := tracer.Start(context.Background(), "request")
	env, err := jobs.Encode(reqCtx, greet{Name: "ada"})
	request.End()
	if err != nil {
		t.Fatal(err)
	}
	if env.Trace["traceparent"] == "" {
		t.Fatalf("envelope carries no trace: %+v", env)
	}

	reg := jobs.NewRegistry()
	jobs.Register(reg, func(context.Context, greet) error { return nil })
	deliveryCtx, delivery := tracer.Start(context.Background(), "delivery")
	err = reg.Dispatch(deliveryCtx, env)
	delivery.End()
	if err != nil {
		t.Fatal(err)
	}

	var job sdktrace.ReadOnlySpan
	for _, s := range rec.Ended() {
		if s.Name() == "job greet" {
			job = s
		}
	}
	if job == nil {
		t.Fatal("no job span")
	}
	if job.Parent().SpanID() != request.SpanContext().SpanID() || job.SpanContext().TraceID() != request.SpanContext().TraceID() {
		t.Errorf("job span's parent = %v, want the request span", job.Parent())
	}
	links := job.Links()
	if len(links) != 1 || links[0].SpanContext.SpanID() != delivery.SpanContext().SpanID() {
		t.Errorf("job span links = %+v, want the delivery span", links)
	}
}
