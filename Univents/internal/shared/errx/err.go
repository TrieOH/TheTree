package errx

import (
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)

type Error struct {
	Kind       string
	Resource   string
	Field      string
	Constraint string
	Message    string
	Cause      error
}

func (e Error) Error() string {
	if e.Message != "" {
		return e.Message
	}

	if e.Constraint != "" {
		return e.Kind + ":" + e.Constraint
	}

	return e.Resource + " " + e.Kind
}

type CombinedError struct {
	errorType string
	errors    []error
}

func (e *CombinedError) Error() string {
	if len(e.errors) == 0 {
		return ""
	}

	var msg strings.Builder
	if e.errorType != "" {
		msg.WriteString("multiple " + e.errorType + " errors:")

	} else {
		msg.WriteString("multiple errors:")
	}

	for _, err := range e.errors {
		msg.WriteString("\n - ")
		msg.WriteString(err.Error())
	}

	return msg.String()
}

func (e *CombinedError) Unwrap() []error {
	return e.errors
}

func Combine(errorType string, errs ...error) error {
	var filtered []error

	for _, err := range errs {
		if err == nil {
			continue
		}

		// Flatten nested CombinedErrors
		var ce *CombinedError
		if errors.As(err, &ce) {
			filtered = append(filtered, ce.errors...)
			continue
		}

		filtered = append(filtered, err)
	}

	if len(filtered) == 0 {
		return nil
	}

	return &CombinedError{
		errorType: errorType,
		errors:    filtered,
	}
}

/*
=== Fluent modifiers ===
*/

func (e Error) SetField(f string) Error {
	e.Field = f
	return e
}

func (e Error) SetConstraint(c string) Error {
	e.Constraint = c
	return e
}

func (e Error) SetMessage(m string) Error {
	e.Message = m
	return e
}

func (e Error) SetCause(err error) Error {
	e.Cause = err
	return e
}

/*
=== Entry means (Kinds)
*/

func Conflict(resource string) Error {
	return Error{
		Kind:     "conflict",
		Resource: resource,
	}
}

func Invalid(resource string) Error {
	return Error{
		Kind:     "invalid",
		Resource: resource,
	}
}

func Forbidden(resource string) Error {
	return Error{
		Kind:     "forbidden",
		Resource: resource,
	}
}

func NotFound(resource string) Error {
	return Error{
		Kind:     "not_found",
		Resource: resource,
	}
}

func Internal(resource string) Error {
	return Error{
		Kind:     "internal",
		Resource: resource,
	}
}

/*
=== DB helpers (pgx-native) ===
*/

func FromDB(err error, resource string) Error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return Internal(resource).SetCause(err)
	}

	switch pgErr.Code {

	// unique_violation
	case "23505":
		return Conflict(resource).
			SetConstraint(pgErr.ConstraintName).
			SetCause(err)

	// foreign_key_violation
	case "23503":
		return Invalid(resource).
			SetConstraint(pgErr.ConstraintName).
			SetCause(err)

	// check_violation
	case "23514":
		return Invalid(resource).
			SetConstraint(pgErr.ConstraintName).
			SetCause(err)

	default:
		return Internal(resource).SetCause(err)
	}
}

/*
=== Specific helpers (optional but nice ergonomics)
You can use these if you already KNOW the context.
*/

func UniqueViolation(resource, constraint string, err error) Error {
	return Conflict(resource).
		SetConstraint(constraint).
		SetCause(err)
}

func ForeignKeyViolation(resource, constraint string, err error) Error {
	return Invalid(resource).
		SetConstraint(constraint).
		SetCause(err)
}

func CheckViolation(resource, constraint string, err error) Error {
	return Invalid(resource).
		SetConstraint(constraint).
		SetCause(err)
}

/*
=== Helpers ===
*/

func As(err error) (Error, bool) {
	var appErr Error
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return Error{}, false
}

func IsKind(err error, kind string) bool {
	var appErr Error
	if errors.As(err, &appErr) {
		return appErr.Kind == kind
	}
	return false
}
