package errx

import (
	"errors"

	fun "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

/*
=== DB helpers (pgx-native) ===
*/

type ConstraintInfo struct {
	Message string
	Type    string // "unique", "fk", "check"
}

var constraintRegistry = map[string]ConstraintInfo{
	// Unique violations
	"idx_forms_title_project": {
		Message: "A form with this title already exists in the project",
		Type:    "unique",
	},
	"uniq_version_number": {
		Message: "Version number already exists for this form",
		Type:    "unique",
	},
	"one_version_draft_per_form": {
		Message: "A draft version already exists",
		Type:    "unique",
	},

	// Check violations
	"chk_forms_status_timestamps": {
		Message: "Invalid status transition: missing required timestamp",
		Type:    "check",
	},
	"chk_fields_key_format": {
		Message: "Field key must start with lowercase letter, contain only letters/numbers/underscores",
		Type:    "check",
	},
	"chk_fields_select_options": {
		Message: "Select fields must have valid options array",
		Type:    "check",
	},
}

func RegisterConstraint(name string, info ConstraintInfo) {
	constraintRegistry[name] = info
}

func FromDB(err error, resource string) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return fun.NewError(resource + " not found").NotFound()
	}

	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return fun.NewError(err.Error()).Internal()
	}

	// Check registry first for custom message
	info, hasCustom := constraintRegistry[pgErr.ConstraintName]

	switch pgErr.Code {
	case "23505": // unique_violation
		msg := pgErr.Error()
		if hasCustom && info.Type == "unique" {
			msg = resource + ": " + info.Message
		}
		return fun.NewError(msg).Conflict()
	case "23503": // foreign_key_violation
		msg := pgErr.Error()
		if hasCustom && info.Type == "fk" {
			msg = resource + ": " + info.Message
		}
		return fun.NewError(msg).BadRequest()
	case "23514": // check_violation
		msg := "Validation failed"
		if hasCustom && info.Type == "check" {
			msg = resource + ": " + info.Message
		}
		return fun.NewError(msg).Validation()
	default:
		return fun.NewError(pgErr.Error()).Internal()
	}
}
