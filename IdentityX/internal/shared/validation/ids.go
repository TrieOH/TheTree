package validation

import (
	"IdentityX/internal/shared/errx"

	"github.com/MintzyG/fail/v3"
	"github.com/google/uuid"
)

// FIXME make these accept ctx and log to span

func ParseRefreshJTI(refreshJTI string) (uuid.UUID, error) {
	jti, err := uuid.Parse(refreshJTI)
	if err != nil {
		return uuid.Nil, fail.New(errx.RequestParseUUIDError).WithArgs("refreshJTI").With(err)
	}
	return jti, nil
}

func RequireRefreshJTI(refreshJTI *string) (uuid.UUID, error) {
	if refreshJTI == nil {
		return uuid.Nil, fail.New(errx.ValidationUUIDWasNil).WithArgs("refreshJTI")
	}
	return ParseRefreshJTI(*refreshJTI)
}

func ParseAccessJTI(accessJTI string) (uuid.UUID, error) {
	jti, err := uuid.Parse(accessJTI)
	if err != nil {
		return uuid.Nil, fail.New(errx.RequestParseUUIDError).WithArgs("accessJTI").With(err)
	}
	return jti, nil
}

func RequireAccessJTI(accessJTI *string) (uuid.UUID, error) {
	if accessJTI == nil {
		return uuid.Nil, fail.New(errx.ValidationUUIDWasNil).WithArgs("accessJTI")
	}
	return ParseAccessJTI(*accessJTI)
}
