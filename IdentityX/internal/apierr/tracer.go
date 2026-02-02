package apierr

import (
	"context"
	"errors"

	"github.com/MintzyG/fail"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type OTelTracer struct {
	Tracer              trace.Tracer
	RecordSystemAsError bool
	RecordDomainAsEvent bool
}

func DefaultOTelTracer() *OTelTracer {
	return &OTelTracer{
		Tracer:              otel.Tracer("fail-example"),
		RecordSystemAsError: true,
		RecordDomainAsEvent: true,
	}
}

func (o *OTelTracer) Trace(operation string, fn func() error) error {
	return o.TraceCtx(context.Background(), operation, func(context.Context) error {
		return fn()
	})
}

func (o *OTelTracer) TraceCtx(ctx context.Context, operation string, fn func(context.Context) error) error {
	ctx, span := o.Tracer.Start(ctx, operation)
	defer span.End()

	err := fn(ctx)
	if err == nil {
		return nil
	}

	o.recordError(span, err)
	return err
}

func (o *OTelTracer) recordError(span trace.Span, err error) {
	if span == nil || err == nil {
		return
	}

	var fe *fail.Error
	if errors.As(err, &fe) {
		if fail.IsDomain(fe) {
			span.SetStatus(codes.Error, fe.Message)
			span.SetAttributes()

			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
			return
		}
		span.AddEvent("domain.error",
			trace.WithAttributes(),
		)

		span.AddEvent("domain.error",
			trace.WithAttributes(
				attribute.String("error.message", err.Error()),
			),
		)
		return

	}
}
