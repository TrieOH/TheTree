package apierr

import (
	"github.com/jackc/pgx/v5/pgconn"
)

func fromUniqueViolation(pgErr *pgconn.PgError, cause error) *Error {
	switch pgErr.ConstraintName {

	case "one_version_draft_per_schema":
		return ErrConflict.
			WithMsg("a draft schema version already exists").
			WithID(ID(SchemaVersionDraftAlreadyExists.String())).
			WithDebugCause(cause)

	case "schema_fields_schema_version_id_position_key":
		return ErrConflict.
			WithMsg("two fields can't occupy the same position").
			WithID(FieldSamePositionForMultipleFields).
			WithDebugCause(cause)

	case "schema_fields_schema_version_id_key_key":
		return ErrConflict.
			WithMsg("two fields can't have the same key").
			WithID(FieldSameKeyForMultipleFields).
			WithDebugCause(cause)

	case "scopes_unique_project_resource_scopes":
		return ErrConflict.
			WithMsg("two scopes can't have the same name and external_id").
			WithID(ScopeDuplicateNameAndExternalID).
			WithDebugCause(cause)

	case "roles_name_project_id_key":
		return ErrConflict.
			WithMsg("role name already taken").
			WithID(RoleNameTaken).
			WithDebugCause(cause)

	case "identity_roles_identity_id_role_id_scope_id_key":
		return ErrConflict.
			WithMsg("user already has this role in the specified scope").
			WithID(RoleAlreadyGranted).
			WithDebugCause(cause)

	case "identity_permissions_identity_id_permission_id_scope_id_key":
		return ErrConflict.
			WithMsg("user already has this permission in the specified scope").
			WithID(PermissionAlreadyGranted).
			WithDebugCause(cause)

	default:
		return ErrConflict.
			WithMsg("resource already exists").
			WithID(DBUniqueViolation).
			WithDebugCause(cause)
	}
}

func fromCheckViolation(pgErr *pgconn.PgError, cause error) *Error {
	switch pgErr.ConstraintName {

	case "schema_fields_key_check":
		return ErrInvalidInput.
			WithMsg("field key must start with a lowercase letter and contain only lowercase letters, numbers, or underscores").
			WithID(FieldInvalidCharactersInKey).
			WithDebugCause(cause)

	case "scope_shape_check":
		return ErrInvalidInput.
			WithMsg("invalid scope shape: a scope must be one of the following — (1) a global scope with type='global' and no project_id, name, or external_id; (2) a project root scope with type='project_root', a project_id, and no name or external_id; or (3) a project scope with type='project_scope', a project_id, and a name (external_id optional)").
			WithID(ScopeInvalid).
			WithDebugCause(cause)

	default:
		return ErrInvalidInput.
			WithMsg("invalid value violates a database constraint").
			WithID(DBCheckViolation).
			WithDebugCause(cause)
	}
}
