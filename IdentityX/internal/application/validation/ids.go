package validation

import (
	"fmt"

	"github.com/google/uuid"
)

// FIXME: Make parse JTI standalone errors

type ErrParseUUID struct {
	FieldName string
	Cause     error
}

func (e ErrParseUUID) Error() string {
	return fmt.Sprintf("invalid uuid field: %s", e.FieldName)
}

func ParseRefreshJTI(refreshJTI string) (uuid.UUID, error) {
	jti, err := uuid.Parse(refreshJTI)
	if err != nil {
		return uuid.Nil, ErrParseUUID{FieldName: "refreshJTI"}
	}
	return jti, nil
}

func RequireRefreshJTI(refreshJTI *string) (uuid.UUID, error) {
	return ParseRefreshJTI(*refreshJTI)
}

func ParseAccessJTI(refreshJTI string) (uuid.UUID, error) {
	jti, err := uuid.Parse(refreshJTI)
	if err != nil {
		return uuid.Nil, ErrParseUUID{FieldName: "accessJTI"}
	}
	return jti, nil
}

func RequireAccessJTI(refreshJTI *string) (uuid.UUID, error) {
	return ParseAccessJTI(*refreshJTI)
}
