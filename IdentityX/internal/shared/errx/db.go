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
		AND tablename IN ('projects', 'users', 'sessions', 'key_pair', 'token_reuse_list', 'api_keys')
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
	// users
	"chk_valid_user_type":   "user type must be one of: client, project",
	"one_email_for_client":  "an account with this email already exists",
	"one_email_per_project": "an account with this email already exists in this project",

	// sessions
	"chk_session_valid_user_type":           "session user type must be one of: client, project",
	"chk_session_not_revoked_before_issued": "a session cannot be revoked before it was issued",
	"sessions_token_id_key":                 "a session with this token ID already exists",

	// key_pair
	"chk_key_pair_key_type_valid":                 "key type must be one of: goauth, project",
	"chk_key_pair_usage_valid":                    "key usage must be one of: sign, verify",
	"chk_key_pair_status_valid":                   "key status must be one of: active, rotated, revoked",
	"chk_key_pair_type_project_consistency_check": "goauth keys must not have a project, project keys must have a project",
	"chk_key_pair_cant_sign_if_rotated":           "a rotated key pair cannot be used for signing",
	"key_pair_kid_key":                            "a key pair with this kid already exists",
	"one_identity_x_active_signing_key":           "there can only be one active goauth signing key",
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
