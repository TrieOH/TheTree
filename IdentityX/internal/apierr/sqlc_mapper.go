package apierr

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func FromSQLC(err error) *Error {
	if err == nil {
		return nil
	}

	// sqlc "not found"
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound.
			WithMsg("resource not found").
			WithID(DBNotFound).
			WithDebugCause(err)
	}

	// postgres error
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return fromPQError(pqErr, err)
	}

	// fallback
	return ErrInternal.
		WithMsg("internal database error").
		WithID(SystemInternalError).
		WithDebugCause(err)
}

func fromPQError(pqErr *pq.Error, cause error) *Error {
	switch string(pqErr.Code) {

	// ---------- constraints ----------
	case "23505": // unique_violation
		return fromUniqueViolation(pqErr, cause)

	case "23514":
		return fromCheckViolation(pqErr, cause)

	case "23503": // foreign_key_violation
		return ErrInvalidInput.
			WithMsg("invalid reference").
			WithID(DBForeignKeyViolation).
			WithDebugCause(cause)

	case "23502": // not_null_violation
		return ErrInvalidInput.
			WithMsg("missing required field").
			WithID(DBNotNullViolation).
			WithDebugCause(cause)

	// ---------- data ----------
	case "22001": // string_data_right_truncation
		return ErrInvalidInput.
			WithMsg("value too long").
			WithID(DBValueTooLong).
			WithDebugCause(cause)

	// ---------- concurrency ----------
	case "40001": // serialization_failure
		return ErrConflict.
			WithMsg("transaction conflict, retry").
			WithID(DBSerializationFailure).
			WithDebugCause(cause)

	// ---------- connection ----------
	case "08006", "08001":
		return ErrInternal.
			WithMsg("database connection error").
			WithID(SystemDependencyDown).
			WithDebugCause(cause)
	}

	// unknown postgres error
	return ErrInternal.
		WithMsg("database error").
		WithID(SystemInternalError).
		WithDebugCause(cause)
}

func RecordSQLCError(span trace.Span, err error) {
	if span == nil || err == nil {
		return
	}

	var apiErr *Error
	if !errors.As(err, &apiErr) {
		// unknown error, treat as system
		RecordSystemError(span, err)
		return
	}

	// Always annotate DB context
	span.SetAttributes(
		attribute.String("error.source", "database"),
		attribute.String("error.code", string(apiErr.Code)),
		attribute.String("error.id", string(apiErr.ID)),
	)

	if IsSystemError(apiErr) {
		span.SetStatus(codes.Error, apiErr.Message)
		for _, subErr := range apiErr.DebugCauses {
			if subErr != nil {
				span.RecordError(subErr)
			}
		}
	}

	// domain-level DB error
	span.AddEvent("domain.error",
		trace.WithAttributes(
			attribute.String("error.code", string(apiErr.Code)),
			attribute.String("error.id", string(apiErr.ID)),
			attribute.String("error.message", apiErr.Message),
		),
	)
}
