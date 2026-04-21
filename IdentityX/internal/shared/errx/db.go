package errx

import (
	"database/sql"
	"errors"

	fun "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/jackc/pgx/v5/pgconn"
)

func FromDB(err error, resource string) error {
	if errors.Is(err, sql.ErrNoRows) {
		return fun.NewErrorf("%s, not found", resource).WithErr(err).NotFound()
	}
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return fun.NewError("unrecognized database error").WithErr(err).Internal()
	}
	switch pgErr.Code {
	case "23505":
		return fun.NewError(pgErr.ConstraintName).WithErr(err).Conflict() // unique_violation
	case "23503":
		return fun.NewError(pgErr.ConstraintName).WithErr(err).BadRequest() // foreign_key_violation
	case "23514":
		return fun.NewError(pgErr.ConstraintName).WithErr(err).BadRequest() // check_violation
	default:
		return fun.NewError("internal database error").WithErr(err).Internal()
	}
}
