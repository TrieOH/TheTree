package errx

import (
	"context"
	"errors"
	"fmt"

	"github.com/MintzyG/FastUtilitiesNet"
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
		return fmt.Errorf("constraints missing from registry: %v", missing)
	}
	return nil
}

var constraintRegistry = map[string]string{
	"chk_forms_status":                "Invalid form status, must be one of: draft, open, closed, archived",
	"chk_forms_valid_status_state":    "Invalid status transition: missing required timestamp for current status",
	"chk_version_gt_zero":             "Version number must be greater than zero",
	"chk_versions_status":             "Invalid version status, must be one of: draft, active, outdated",
	"chk_fields_select_options":       "Select fields must have a valid non-empty options array",
	"chk_fields_key_format":           "Field key must start with a lowercase letter and contain only letters, numbers, and underscores",
	"chk_fields_type":                 "Invalid field type, must be one of: string, email, int, float, bool, select",
	"chk_fields_owner":                "Invalid field owner, must be one of: user, admin",
	"chk_fields_select_behaviour":     "Select fields must specify a behaviour: checkbox or radio",
	"uniq_one_key_per_version":        "A field with this key already exists in this version",
	"uniq_one_stable_per_version":     "A field with this stable ID already exists in this version",
	"one_version_active_per_form":     "An active version already exists for this form",
	"one_version_draft_per_form":      "A draft version already exists for this form",
	"uniq_idx_api_keys_name_project":  "An active API key with this name already exists in the project",
	"uniq_idx_forms_title_project":    "A form with this title already exists in the project",
	"uniq_idx_projects_owner_id_name": "A project with this name already exists",
	"uniq_idx_version_number":         "This version number already exists for this form",
}

func DB(err error, resource string) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return fun.NewError(resource + " not found").NotFound()
	}

	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return fun.NewError(err.Error()).Internal()
	}

	constraintMessage, hasCustom := constraintRegistry[pgErr.ConstraintName]
	switch pgErr.Code {
	case "23505": // unique_violation
		msg := resource + ": " + pgErr.Error()
		if hasCustom {
			msg = constraintMessage
		}
		return fun.NewError(msg).Conflict()
	case "23503": // foreign_key_violation
		msg := resource + ": " + pgErr.Error()
		if hasCustom {
			msg = constraintMessage
		}
		return fun.NewError(msg).BadRequest()
	case "23514": // check_violation
		msg := resource + ": Validation failed"
		if hasCustom {
			msg = constraintMessage
		}
		return fun.NewError(msg).Validation()
	default:
		return fun.NewError(pgErr.Error()).Internal()
	}
}
