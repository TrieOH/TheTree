package database

import (
	"errors"

	"github.com/MintzyG/fun"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

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

		constraintMessage, hasCustom := ConstraintErrorRegistry[pgErr.ConstraintName]
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
