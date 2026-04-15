package errx

import (
	"IdentityX/internal/platform/telemetry"
	"database/sql"
	"errors"

	"github.com/MintzyG/fail/v3"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

type PGXMapper struct {
	priority int
	name     string
}

func (m *PGXMapper) Name() string  { return m.name }
func (m *PGXMapper) Priority() int { return m.priority }

func (m *PGXMapper) Map(err error) (fe *fail.Error, ok bool) {
	if IsNotFound(err) {
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
			case "users_email_key":
				return fail.New(AuthEmailAlreadyUsed).Msg("email already in use").Debug(err.Error()), true
			case "project_users_project_id_email_key":
				return fail.New(AuthEmailAlreadyUsed).Msg("email already in use").Debug(err.Error()), true
			case "permissions_project_id_object_action_key":
				return fail.New(PERMissionAlreadyExists).Debug(err.Error()), true
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
				return fail.New(PERMissionAlreadyGranted).Debug(err.Error()), true
			case "role_permissions_pkey":
				return fail.New(ROLEPermissionAlreadyGranted).Debug(err.Error()), true
			case "scopes_one_global":
				return fail.New(SCOPEOneGlobal).Debug(err.Error()), true
			case "scopes_one_project_root_per_project":
				return fail.New(SCOPEOneProjectRootPerProject).Debug(err.Error()), true
			case "scopes_unique_siblings", "scopes_unique_resource_siblings":
				return fail.New(SCOPEDuplicateSibling).Debug(err.Error()), true
			case "uniq_goauth_single_active_signing_key":
				return fail.New(SCOPEDuplicateSibling).Debug(err.Error()), true
			default:
				telemetry.Log().Info("error", zap.Error(err))
				panic(err.Error())
				//return fail.New(SQLUnmatchedUniqueViolation).Debug(err.Error()), true
			}
		case "23514": // check_violation
			defer func() {
				_ = fe.AddMeta("is_check_violation", true)
			}()
			switch pgErr.ConstraintName {
			case "schema_fields_key_check":
				return fail.New(FIELDInvalidCharactersInKey).Debug(err.Error()), true
			case "scope_shape_check":
				return fail.New(SCOPEInvalidShape).Debug(err.Error()), true
			default:
				return fail.New(SQLUnmatchedCheckViolation).Debug(err.Error()), true
			}
		case "23503": // foreign key violation
			switch pgErr.ConstraintName {
			case "project_users_project_id_fkey":
				return fail.New(ProjectUserRegisterOnNoneProject).Debug(err.Error()), true
			}
			return fail.New(SQLForeignKeyViolation).With(err).Debug(err.Error()), true
		case "23502":
			return fail.New(SQLNotNULLViolation).With(err).Debug(err.Error()), true
		case "22001":
			return fail.New(SQLValueTooLong).With(err).Debug(err.Error()), true // data
		case "40001":
			return fail.New(SQLSerializationFailure).With(err).Debug(err.Error()), true // concurrency
		case "08006", "08001", "08004":
			return fail.New(SQLDBConnectionError).With(err).Debug(err.Error()), true // connection
		}

		return fail.New(SQLUnknownError).With(err).Debug(err.Error()), true
	}

	return nil, false
}

func IsNotFound(err error) bool {
	if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
		return true
	}
	var fe *fail.Error
	if errors.As(err, &fe) {
		if fail.Is(fe, SQLNotFound) {
			return true
		}
	}
	return false
}

func IsUniqueViolation(err error) bool {
	if err != nil {
		return false
	}
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
	if err != nil {
		return false
	}
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
