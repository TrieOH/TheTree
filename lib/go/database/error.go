package database

import (
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// The database errors every Backend shares (ADR-0016). An ErrorHandler
// joins them under the Backend's named errors, so callers can match either
// the specific error or its class.
var (
	ErrNotFound            = errors.New("not found")
	ErrUniqueViolation     = errors.New("conflicts with an existing record")
	ErrForeignKeyViolation = errors.New("references a record that does not exist")
	ErrCheckViolation      = errors.New("invalid value")
	ErrExclusionViolation  = errors.New("conflicts with an existing record")
)

// ConstraintError is a violation of a constraint the registry has no named
// error for. It matches its class (errors.Is) and keeps the driver error for
// logs; its constraint name and the driver text never reach a client.
type ConstraintError struct {
	Class      error
	Constraint string
	Err        *pgconn.PgError
}

func (e *ConstraintError) Error() string {
	return e.Class.Error() + " (" + e.Constraint + ")"
}

func (e *ConstraintError) Unwrap() []error { return []error{e.Class, e.Err} }

// NotFoundError is a missing row: it reads as the resource's not-found
// error and matches both it and ErrNotFound (errors.Is).
type NotFoundError struct{ Named error }

func (e *NotFoundError) Error() string { return e.Named.Error() }

func (e *NotFoundError) Unwrap() []error { return []error{e.Named, ErrNotFound} }

// ViolationError is a registered constraint's violation: it reads as the named
// error and matches both it and its class (errors.Is), keeping the driver
// error for logs.
type ViolationError struct {
	Named error
	Class error
	Err   *pgconn.PgError
}

func (e *ViolationError) Error() string { return e.Named.Error() }

func (e *ViolationError) Unwrap() []error { return []error{e.Named, e.Class, e.Err} }

// ErrorHandler turns a driver error into the Backend's named errors
// (ADR-0016). Repos hold one per resource and return repo.dbe(err).
type ErrorHandler func(err error, notFound ...error) error

// NewErrorHandler builds the handler for one resource. notFound is the
// resource's named not-found error; a call can override it (a lookup of a
// different resource through the same repo).
//
//   - pgx.ErrNoRows: a NotFoundError (the not-found error and ErrNotFound).
//   - A registered constraint: a ViolationError (its named error and class).
//   - An unregistered constraint: a ConstraintError of its class.
//   - Anything else: returned as is (an internal error).
func NewErrorHandler(notFound error) ErrorHandler {
	return func(err error, override ...error) error {
		if err == nil {
			return nil
		}
		if errors.Is(err, pgx.ErrNoRows) {
			nf := notFound
			if len(override) > 0 {
				nf = override[0]
			}
			return &NotFoundError{Named: nf}
		}
		pgErr, ok := errors.AsType[*pgconn.PgError](err)
		if !ok {
			return err
		}
		class := constraintClass(pgErr.Code)
		named, registered := constraintRegistry[pgErr.ConstraintName]
		switch {
		case class == nil:
			return err
		case registered:
			return &ViolationError{Named: named, Class: class, Err: pgErr}
		default:
			return &ConstraintError{Class: class, Constraint: pgErr.ConstraintName, Err: pgErr}
		}
	}
}

func constraintClass(code string) error {
	switch code {
	case "23505":
		return ErrUniqueViolation
	case "23503":
		return ErrForeignKeyViolation
	case "23514":
		return ErrCheckViolation
	case "23P01":
		return ErrExclusionViolation
	}
	return nil
}

// BatchError is a batch with failed items, each named by its id.
type BatchError struct{ Items []BatchItemError }

// BatchItemError is one failed item of a batch.
type BatchItemError struct {
	ID  string
	Err error
}

func (e *BatchError) Error() string {
	msgs := make([]string, len(e.Items))
	for i, item := range e.Items {
		msgs[i] = item.ID + ": " + item.Err.Error()
	}
	return strings.Join(msgs, "; ")
}

func (e *BatchError) Unwrap() []error {
	errs := make([]error, len(e.Items))
	for i, item := range e.Items {
		errs[i] = item.Err
	}
	return errs
}

// BatchExec runs a batch and collects every failed item, mapped by
// dbe, into one BatchError (nil when all items succeed).
func BatchExec(br interface {
	Exec(func(int, error))
	Close() error
}, dbe ErrorHandler, getID func(int) string) error {
	defer func() {
		_ = br.Close()
	}()
	var items []BatchItemError
	br.Exec(func(i int, err error) {
		if err != nil {
			items = append(items, BatchItemError{ID: getID(i), Err: dbe(err)})
		}
	})
	if len(items) == 0 {
		return nil
	}
	return &BatchError{Items: items}
}
