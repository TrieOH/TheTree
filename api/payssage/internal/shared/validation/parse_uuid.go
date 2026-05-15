package validation

import (
	"payssage/internal/shared/errx"

	"github.com/google/uuid"
)

func ParseUUID(id, fieldName string) (uuid.UUID, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, errx.Invalid("uuid").SetCause(err)
	}
	return parsedID, nil
}
