package apierr

import (
	"errors"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Code string
type ID string

type Error struct {
	Code    Code
	ID      ID
	Message string
	Fields  map[string]any
	Cause   error
}

func (e Error) Error() string {
	if e.Cause == nil {
		return string(e.ID) + ": " + e.Message
	}
	return string(e.ID) + ": " + e.Message + "; CAUSE (" + e.Cause.Error() + ") "
}

func (e Error) Unwrap() error {
	return e.Cause
}

func RecordSystemError(span trace.Span, err error) {
	if span == nil || err == nil {
		return
	}

	var apiErr *Error
	if errors.As(err, &apiErr) {
		span.SetStatus(codes.Error, apiErr.Message)
		span.SetAttributes(
			attribute.String("error.code", string(apiErr.Code)),
			attribute.String("error.id", string(apiErr.ID)),
		)

		if apiErr.Cause != nil {
			span.RecordError(apiErr.Cause)
			return
		}
	}

	span.SetStatus(codes.Error, err.Error())
	span.RecordError(err)
}

func RecordDomainError(span trace.Span, err error) {
	if span == nil || err == nil {
		return
	}

	var apiErr *Error
	if errors.As(err, &apiErr) {
		span.AddEvent("domain.error",
			trace.WithAttributes(
				attribute.String("error.code", string(apiErr.Code)),
				attribute.String("error.id", string(apiErr.ID)),
				attribute.String("error.message", apiErr.Message),
			),
		)
		return
	}

	span.AddEvent("domain.error",
		trace.WithAttributes(
			attribute.String("error.message", err.Error()),
		),
	)
}
