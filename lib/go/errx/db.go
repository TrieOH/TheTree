package errx

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/MintzyG/fun"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ConstraintRegistry map[string]string

type DBHandler struct {
	registry ConstraintRegistry
	tables   []string
}

func NewDBHandler(registry ConstraintRegistry, tables []string) *DBHandler {
	return &DBHandler{registry: registry, tables: tables}
}

func (h *DBHandler) Validate(ctx context.Context, db *pgxpool.Pool) error {
	tableList := "'" + strings.Join(h.tables, "', '") + "'"
	rows, err := db.Query(ctx, fmt.Sprintf(`
        SELECT con.conname
        FROM pg_constraint con
        JOIN pg_class rel ON rel.oid = con.conrelid
        JOIN pg_namespace nsp ON nsp.oid = rel.relnamespace
        WHERE nsp.nspname = 'public'
        AND con.contype IN ('u', 'c')
        UNION
        SELECT indexname
        FROM pg_indexes
        WHERE schemaname = 'public'
        AND tablename IN (%s)
        AND indexname LIKE 'uniq_%%' OR indexname LIKE 'one_%%'
    `, tableList))
	if err != nil {
		return err
	}
	defer rows.Close()

	var missing []string
	for rows.Next() {
		var name string
		_ = rows.Scan(&name)
		if _, ok := h.registry[name]; !ok {
			missing = append(missing, name)
		}
	}

	if len(missing) > 0 {
		return fun.Errf("constraints missing from registry: %v", missing).Internal()
	}
	return nil
}

func (h *DBHandler) DB(err error, resource string) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return fun.ErrNotFound(resource + " not found")
	}

	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return fun.ErrInternal(err.Error())
	}

	constraintMessage, hasCustom := h.registry[pgErr.ConstraintName]
	switch pgErr.Code {
	case "23505":
		msg := resource + ": " + pgErr.Error()
		if hasCustom {
			msg = constraintMessage
		}
		return fun.ErrConflict(msg)
	case "23503":
		msg := resource + ": " + pgErr.Error()
		if hasCustom {
			msg = constraintMessage
		}
		return fun.ErrBadRequest(msg)
	case "23514":
		msg := resource + ": Validation failed"
		if hasCustom {
			msg = constraintMessage
		}
		return fun.ErrValidation(msg)
	default:
		return fun.ErrInternal(pgErr.Error())
	}
}
