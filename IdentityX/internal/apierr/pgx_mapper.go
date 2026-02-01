package apierr

import (
	"database/sql"
	"errors"

	"github.com/MintzyG/fail"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// FIXME later change fail to check if its static and not let you change anything in the error

type PGXMapper struct {
	priority int
	name     string
}

func (m *PGXMapper) Name() string  { return m.name }
func (m *PGXMapper) Priority() int { return m.priority }

func (m *PGXMapper) Map(err error) (error, bool) {
	return m.MapToFail(err)
}

func (m *PGXMapper) MapFromFail(fe *fail.Error) (error, bool) {
	return errors.New(fe.Message), true
}

func (m *PGXMapper) MapToFail(err error) (fe *fail.Error, ok bool) {
	if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
		return fail.New(SQLNotFound), true
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			defer func() {
				_ = fe.AddMeta("is_unique_violation", true)
			}()
			switch pgErr.ConstraintName {
			case "one_version_draft_per_schema":
				return fail.New(SCHEMAVersionDraftAlreadyExists).Debug(err.Error()), true
			case "schema_fields_schema_version_id_position_key":
				return fail.New(FIELDSamePositionForMultipleFields).Debug(err.Error()), true
			case "schema_fields_schema_version_id_key_key":
				return fail.New(FIELDSameKeyForMultipleFields).Debug(err.Error()), true
			case "scopes_unique_project_resource_scopes":
				return fail.New(SCOPEDuplicateNameAndExternalID).Debug(err.Error()), true
			case "roles_name_project_id_key":
				return fail.New(ROLENameAlreadyTaken).Debug(err.Error()), true
			case "identity_roles_identity_id_role_id_scope_id_key":
				return fail.New(ROLEAlreadyGranted).Debug(err.Error()), true
			case "identity_permissions_identity_id_permission_id_scope_id_key":
				return fail.New(PERMISSIONAlreadyGranted).Debug(err.Error()), true
			default:
				return fail.New(SQLUnmatchedUniqueViolation).Debug(err.Error()), true
			}
		case "23514": // check_violation
			defer func() {
				_ = fe.AddMeta("is_check_violation", true)
			}()
			switch pgErr.ConstraintName {
			case "schema_fields_key_check":
				return fail.New(FIELDInvalidCharactersInKey).Debug(err.Error()), true
			case "scope_shape_check":
				return fail.New(SCOPEInvalid).Debug(err.Error()), true
			default:
				return fail.New(SQLUnmatchedCheckViolation).Debug(err.Error()), true
			}
		case "23503":
			return fail.New(SQLForeignKeyViolation).With(err), true
		case "23502":
			return fail.New(SQLNotNULLViolation).With(err), true
		case "22001":
			return fail.New(SQLValueTooLong).With(err), true // data
		case "40001":
			return fail.New(SQLSerializationFailure).With(err), true // concurrency
		case "08006", "08001", "08004":
			return fail.New(SQLDBConnectionError).With(err), true // connection
		}

		return fail.New(SQLUnknownError).With(err), true
	}

	return nil, false
}

func IsUniqueViolationNew(err error) bool {
	var fe *fail.Error
	if errors.As(err, &fe) {
		if fe.Meta != nil {
			isUnique, ok := fe.Meta["is_unique_violation"]
			if ok {
				return isUnique.(bool)
			}
		}
	}
	return false
}

func IsCheckViolationNew(err error) bool {
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
