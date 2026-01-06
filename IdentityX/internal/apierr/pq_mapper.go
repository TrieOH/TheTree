package apierr

import "github.com/lib/pq"

func fromUniqueViolation(pqErr *pq.Error, cause error) *Error {
	switch pqErr.Constraint {

	case "one_version_draft_per_schema":
		return ErrConflict.
			WithMsg("a draft schema version already exists").
			WithID(SchemaVersionDraftAlreadyExists).
			WithCause(cause)

	default:
		return ErrConflict.
			WithMsg("resource already exists").
			WithID(DBUniqueViolation).
			WithCause(cause)
	}
}
