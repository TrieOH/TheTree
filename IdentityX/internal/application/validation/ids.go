package validation

import (
	"GoAuth/internal/apierr"

	"github.com/MintzyG/fail"
	"github.com/google/uuid"
)

func ParseRefreshJTI(refreshJTI string) (uuid.UUID, error) {
	jti, err := uuid.Parse(refreshJTI)
	if err != nil {
		return uuid.Nil, fail.New(apierr.RequestParseUUIDError).WithArgs("refreshJTI").With(err)
	}
	return jti, nil
}

func RequireRefreshJTI(refreshJTI *string) (uuid.UUID, error) {
	if refreshJTI == nil {
		return uuid.Nil, fail.New(apierr.ValidationUUIDWasNil).WithArgs("refreshJTI")
	}
	return ParseRefreshJTI(*refreshJTI)
}

func ParseAccessJTI(accessJTI string) (uuid.UUID, error) {
	jti, err := uuid.Parse(accessJTI)
	if err != nil {
		return uuid.Nil, fail.New(apierr.RequestParseUUIDError).WithArgs("accessJTI").With(err)
	}
	return jti, nil
}

func RequireAccessJTI(accessJTI *string) (uuid.UUID, error) {
	if accessJTI == nil {
		return uuid.Nil, fail.New(apierr.ValidationUUIDWasNil).WithArgs("accessJTI")
	}
	return ParseAccessJTI(*accessJTI)
}
