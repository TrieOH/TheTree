package apierr

import (
	"go.opentelemetry.io/otel/trace"
)

// FIXME remove me
func FromService(_ trace.Span, err error) error {
	return err
}
