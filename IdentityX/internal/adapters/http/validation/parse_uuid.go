package validation

import (
	"GoAuth/internal/errx"

	"github.com/MintzyG/fail/v3"
	"github.com/google/uuid"
)

func ParseUUID(id, fieldName string) (uuid.UUID, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, fail.New(errx.RequestParseUUIDError).WithArgs(fieldName).With(err)
	}
	return parsedID, nil
}
