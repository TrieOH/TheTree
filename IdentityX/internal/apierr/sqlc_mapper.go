package apierr

import (
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func IsUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}

func IsCheckViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23514"
	}
	return false
}

func FromSQLC(err error) *Error {
	if err == nil {
		return nil
	}

	// sqlc "not found" - this stays the same
	if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound.
			WithMsg("resource not found").
			WithID(DBNotFound).
			WithDebugCause(err)
	}

	// pgx/postgres error
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return fromPGError(pgErr, err)
	}

	// fallback
	return ErrInternal.
		WithMsg("internal database error").
		WithID(PlaceholderID).
		WithDebugCause(err)
}

// renamed from fromPQError to fromPGError
func fromPGError(pgErr *pgconn.PgError, cause error) *Error {
	switch pgErr.Code {

	// ---------- constraints ----------
	case "23505": // unique_violation
		return fromUniqueViolation(pgErr, cause)

	case "23514":
		return fromCheckViolation(pgErr, cause)

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
	case "08006", "08001", "08004":
		return ErrInternal.
			WithMsg("database connection error").
			WithDebugCause(cause)
	}

	// unknown postgres error
	return ErrInternal.
		WithMsg("database error").
		WithID(PlaceholderID).
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
