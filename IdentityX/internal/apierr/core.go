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
	Code        Code
	ID          ID
	Message     string
	Fields      map[string]any
	Cause       error   // Deprecated: Use Causes instead. Kept for backwards compatibility.
	Causes      []error // Multiple causes for error chain
	DebugCauses []error // Debug-only causes (SQL errors, etc.) - safe to disable in prod
}

func (e Error) Error() string {
	if e.Cause == nil && len(e.Causes) == 0 {
		return string(e.ID) + ": " + e.Message
	}

	// Backwards compatibility: use Cause if Causes is empty
	if len(e.Causes) == 0 && e.Cause != nil {
		return string(e.ID) + ": " + e.Message + "; CAUSE (" + e.Cause.Error() + ") "
	}

	// Multiple causes
	msg := string(e.ID) + ": " + e.Message
	for i, cause := range e.Causes {
		msg += "; CAUSE[" + string(rune('0'+i)) + "] (" + cause.Error() + ")"
	}
	return msg
}

func (e Error) Unwrap() error {
	// Backwards compatibility: return Cause if Causes is empty
	if len(e.Causes) == 0 {
		return e.Cause
	}

	// Return first cause for errors.Is/As compatibility
	if len(e.Causes) > 0 {
		return e.Causes[0]
	}

	return nil
}

// GetAllCauses returns all causes including the legacy Cause field
func (e Error) GetAllCauses() []error {
	causes := make([]error, 0, len(e.Causes)+1)

	// Include legacy Cause if it exists and isn't in Causes
	if e.Cause != nil {
		found := false
		for _, c := range e.Causes {
			if c == e.Cause {
				found = true
				break
			}
		}
		if !found {
			causes = append(causes, e.Cause)
		}
	}

	causes = append(causes, e.Causes...)
	return causes
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

		// Record all causes
		allCauses := apiErr.GetAllCauses()
		if len(allCauses) > 0 {
			for _, cause := range allCauses {
				span.RecordError(cause)
			}
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
