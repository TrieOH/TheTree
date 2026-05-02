package errx

import (
	"context"
	"errors"

	"github.com/MintzyG/fun"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ValidateConstraintRegistry(ctx context.Context, db *pgxpool.Pool) error {
	rows, err := db.Query(ctx, `
        SELECT con.conname
        FROM pg_constraint con
        JOIN pg_class rel ON rel.oid = con.conrelid
        JOIN pg_namespace nsp ON nsp.oid = rel.relnamespace
        WHERE nsp.nspname = 'public'
        AND con.contype IN ('u', 'c') -- unique, check
		UNION
		SELECT indexname
		FROM pg_indexes
		WHERE schemaname = 'public'
		AND tablename IN ('projects', 'api_keys', 'forms', 'versions', 'fields')
		AND indexname LIKE 'uniq_%' OR indexname LIKE 'one_%'
    `)
	if err != nil {
		return err
	}
	defer rows.Close()

	var missing []string
	for rows.Next() {
		var name string
		_ = rows.Scan(&name)
		if _, ok := constraintRegistry[name]; !ok {
			missing = append(missing, name)
		}
	}

	if len(missing) > 0 {
		return fun.Errf("constraints missing from registry: %v", missing).Internal()
	}
	return nil
}

var constraintRegistry = map[string]string{
	// forms
	"chk_forms_valid_status":       "status must be one of: draft, open, closed, archived",
	"chk_forms_valid_status_state": "opened_at, closed_at or archived_at must be set when status is open, closed or archived",
	"uniq_form_name_per_namespace": "a form with this name already exists in this namespace",
	"uniq_name_per_user":           "an API key with this name already exists",

	// fields
	"chk_fields_type":             "field type must be one of: string, email, int, float, bool, select",
	"chk_fields_key_format":       "field key must start with a letter or underscore and contain only lowercase letters, digits and underscores",
	"chk_fields_select_behaviour": "select fields must have a behaviour of checkbox or radio",
	"chk_fields_select_options":   "select fields must have a non-empty options array",
	"chk_select_type":             "select_type must be one of: string, email, int, float, bool, select",
	"uniq_key_per_step":           "a field with this key already exists in this step",

	// namespaces
	"uniq_namespace_name_per_user": "a namespace with this name already exists",
}

func DB(err error, resource string) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return fun.ErrNotFound(resource + " not found")
	}

	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return fun.ErrInternal(err.Error())
	}

	constraintMessage, hasCustom := constraintRegistry[pgErr.ConstraintName]
	switch pgErr.Code {
	case "23505": // unique_violation
		msg := resource + ": " + pgErr.Error()
		if hasCustom {
			msg = constraintMessage
		}
		return fun.ErrConflict(msg)
	case "23503": // foreign_key_violation
		msg := resource + ": " + pgErr.Error()
		if hasCustom {
			msg = constraintMessage
		}
		return fun.ErrBadRequest(msg)
	case "23514": // check_violation
		msg := resource + ": Validation failed"
		if hasCustom {
			msg = constraintMessage
		}
		return fun.ErrValidation(msg)
	default:
		return fun.ErrInternal(pgErr.Error())
	}
}
