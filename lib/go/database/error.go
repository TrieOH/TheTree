package database

import (
	"errors"
	"fmt"
	"strings"

	"github.com/MintzyG/fun"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// The database errors every Backend shares (ADR-0016). A NamedErrorHandler
// joins them under the Backend's named errors, so callers can match either
// the specific error or its class.
var (
	ErrNotFound            = errors.New("not found")
	ErrUniqueViolation     = errors.New("conflicts with an existing record")
	ErrForeignKeyViolation = errors.New("references a record that does not exist")
	ErrCheckViolation      = errors.New("invalid value")
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

// NamedErrorHandler turns a driver error into the Backend's named errors
// (ADR-0016). Repos hold one per resource and return repo.dbe(err).
type NamedErrorHandler func(err error, notFound ...error) error

// NewNamedErrorHandler builds the handler for one resource. notFound is the
// resource's named not-found error; a call can override it (a lookup of a
// different resource through the same repo).
//
//   - pgx.ErrNoRows: the not-found error, joined with ErrNotFound.
//   - A registered constraint: its named error, joined with its class.
//   - An unregistered constraint: a ConstraintError of its class.
//   - Anything else: returned as is (an internal error).
func NewNamedErrorHandler(notFound error) NamedErrorHandler {
	return func(err error, override ...error) error {
		if err == nil {
			return nil
		}
		if errors.Is(err, pgx.ErrNoRows) {
			nf := notFound
			if len(override) > 0 {
				nf = override[0]
			}
			return fmt.Errorf("%w: %w", nf, ErrNotFound)
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
			return fmt.Errorf("%w: %w", named, class)
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
	}
	return nil
}

// ErrorHandler is the root-fun form of the handler, for Backends that have
// not moved to named errors yet (ADR-0016).
type ErrorHandler func(error, ...string) error

func NewErrorHandler(resource string) ErrorHandler {
	return func(err error, override ...string) error {
		res := resource
		if len(override) > 0 {
			res = override[0]
		}
		if err == nil {
			return nil
		}

		if errors.Is(err, pgx.ErrNoRows) {
			return fun.ErrNotFound(res + " not found")
		}

		var pgErr *pgconn.PgError
		if !errors.As(err, &pgErr) {
			return fun.ErrInternal(err.Error())
		}

		constraintMessage, hasCustom := "", false
		if named, ok := constraintRegistry[pgErr.ConstraintName]; ok {
			constraintMessage, hasCustom = named.Error(), true
		}
		switch pgErr.Code {
		case "23505":
			msg := res + ": " + pgErr.Error()
			if hasCustom {
				msg = constraintMessage
			}
			return fun.ErrConflict(msg)
		case "23503":
			msg := res + ": " + pgErr.Error()
			if hasCustom {
				msg = constraintMessage
			}
			return fun.ErrBadRequest(msg)
		case "23514":
			msg := res + ": Validation failed"
			if hasCustom {
				msg = constraintMessage
			}
			return fun.ErrValidation(msg)
		default:
			return fun.ErrInternal(pgErr.Error())
		}
	}
}

func BatchExec(br interface {
	Exec(func(int, error))
	Close() error
}, dbe ErrorHandler, getID func(int) string) error {
	defer func() {
		_ = br.Close()
	}()
	var msgs []string
	br.Exec(func(i int, err error) {
		if err != nil {
			msgs = append(msgs, fmt.Sprintf("%s: %s", getID(i), dbe(err).Error()))
		}
	})
	if len(msgs) == 0 {
		return nil
	}
	return fun.ErrValidation(strings.Join(msgs, "; "))
}
