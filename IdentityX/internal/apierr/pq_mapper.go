package apierr

import "github.com/lib/pq"

func fromUniqueViolation(pqErr *pq.Error, cause error) *Error {
	switch pqErr.Constraint {

	case "one_version_draft_per_schema":
		return ErrConflict.
			WithMsg("a draft schema version already exists").
			WithID(SchemaVersionDraftAlreadyExists).
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

	default:
		return ErrConflict.
			WithMsg("resource already exists").
			WithID(DBUniqueViolation).
			WithDebugCause(cause)
	}
}

func fromCheckViolation(pqErr *pq.Error, cause error) *Error {
	switch pqErr.Constraint {

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
