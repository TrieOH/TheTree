package validation

import (
	fun "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/google/uuid"
)

func ParseUUID(id, fieldName string) (uuid.UUID, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, fun.NewErrorf("error parsing %s as uuid", fieldName).Validation()
	}
	return parsedID, nil
}
