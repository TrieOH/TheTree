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

type ErrUUIDWasNil struct {
	FieldName string
}

func (e ErrUUIDWasNil) Error() string {
	return e.FieldName + " field is nil"
}

func ParseRefreshJTI(refreshJTI string) (uuid.UUID, error) {
	jti, err := uuid.Parse(refreshJTI)
	if err != nil {
		return uuid.Nil, ErrParseUUID{FieldName: "refreshJTI", Cause: err}
	}
	return jti, nil
}

func RequireRefreshJTI(refreshJTI *string) (uuid.UUID, error) {
	if refreshJTI == nil {
		return uuid.Nil, ErrUUIDWasNil{FieldName: "refreshJTI"}
	}
	return ParseRefreshJTI(*refreshJTI)
}

func ParseAccessJTI(accessJTI string) (uuid.UUID, error) {
	jti, err := uuid.Parse(accessJTI)
	if err != nil {
		return uuid.Nil, ErrParseUUID{FieldName: "accessJTI", Cause: err}
	}
	return jti, nil
}

func RequireAccessJTI(accessJTI *string) (uuid.UUID, error) {
	if accessJTI == nil {
		return uuid.Nil, ErrUUIDWasNil{FieldName: "accessJTI"}
	}
	return ParseAccessJTI(*accessJTI)
}
