package apierr

import "github.com/lib/pq"

func fromUniqueViolation(pqErr *pq.Error, cause error) *Error {
	switch pqErr.Constraint {

	case "one_version_draft_per_schema":
		return ErrConflict.
			WithMsg("a draft schema version already exists").
			WithID(SchemaVersionDraftAlreadyExists).
			WithCause(cause)

	case "schema_fields_schema_version_id_position_key":
		return ErrConflict.
			WithMsg("two fields can't occupy the same position").
			WithID(FieldSamePositionForMultipleFields).
			WithCause(cause)

	default:
		return ErrConflict.
			WithMsg("resource already exists").
			WithID(DBUniqueViolation).
			WithCause(cause)
	}
}

func fromCheckViolation(pqErr *pq.Error, cause error) *Error {
	switch pqErr.Constraint {

	case "schema_fields_key_check":
		return ErrInvalidInput.
			WithMsg("field key must start with a lowercase letter and contain only lowercase letters, numbers, or underscores").
			WithID(FieldInvalidCharactersInKey).
			WithCause(cause)

	default:
		return ErrInvalidInput.
			WithMsg("invalid value violates a database constraint").
			WithID(DBCheckViolation).
			WithCause(cause)
	}
}
