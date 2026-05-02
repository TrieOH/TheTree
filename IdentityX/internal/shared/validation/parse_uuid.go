package validation

import (
	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func ParseUUID(id, fieldName string) (uuid.UUID, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, fun.Errf("invalid uuid in field (%s): %s", fieldName, err).BadRequest()
	}
	return parsedID, nil
}
