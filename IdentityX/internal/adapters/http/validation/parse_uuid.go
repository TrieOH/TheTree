package validation

import (
	"GoAuth/internal/application/validation"

	"github.com/google/uuid"
)

func ParseUUID(id, fieldName string) (uuid.UUID, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, validation.ErrParseUUID{FieldName: fieldName, Cause: err}
	}
	return parsedID, nil
}
