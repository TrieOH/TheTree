package errx

import (
	"database/sql"
	"errors"

	"github.com/MintzyG/fail/v3"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type PGXMapper struct {
	priority int
	name     string
}

func (m *PGXMapper) Name() string  { return m.name }
func (m *PGXMapper) Priority() int { return m.priority }

func (m *PGXMapper) Map(err error) (fe *fail.Error, ok bool) {
	if IsNotFound(err) {
		return fail.New(SQLResourceNotFound), true
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			defer func() {
				_ = fe.AddMeta("is_unique_violation", true)
			}()
			switch pgErr.ConstraintName {
			default:
				panic(err.Error())
				//return fail.New(SQLUnmatchedUniqueViolation).Debug(err.Error()), true
			}
		case "23514": // check_violation
			defer func() {
				_ = fe.AddMeta("is_check_violation", true)
			}()
			switch pgErr.ConstraintName {
			default:
				return fail.New(SQLUnmatchedCheckViolation).Debug(err.Error()), true
			}
		case "23503": // foreign key violation
			switch pgErr.ConstraintName {
			default:
				return fail.New(SQLForeignKeyViolation).With(err).Debug(err.Error()), true
			}
		case "23502":
			return fail.New(SQLNotNULLViolation).With(err).Debug(err.Error()), true
		case "22001":
			return fail.New(SQLValueTooLong).With(err).Debug(err.Error()), true // data
		case "40001":
			return fail.New(SQLSerializationFailure).With(err).Debug(err.Error()), true // concurrency
		case "08006", "08001", "08004":
			return fail.New(SQLDBConnectionError).With(err).Debug(err.Error()), true // connection
		}

		return fail.New(SQLDatabaseUnknownError).With(err).Debug(err.Error()), true
	}

	return nil, false
}

func IsNotFound(err error) bool {
	if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
		return true
	}
	var fe *fail.Error
	if errors.As(err, &fe) {
		if fail.Is(fe, SQLResourceNotFound) {
			return true
		}
	}
	return false
}

func IsUniqueViolation(err error) bool {
	var fe *fail.Error
	if errors.As(err, &fe) {
		if fe.Meta != nil {
			isUnique, ok := fe.Meta["is_unique_violation"]
			if ok {
				return isUnique.(bool)
			}
			isConflict, ok := fe.Meta["code"].(int)
			if ok && isConflict == 409 {
				return true
			}
		}
	}
	return false
}

func IsCheckViolation(err error) bool {
	var fe *fail.Error
	if errors.As(err, &fe) {
		if fe.Meta != nil {
			isCheck, ok := fe.Meta["is_check_violation"]
			if ok {
				return isCheck.(bool)
			}
		}
	}
	return false
}
