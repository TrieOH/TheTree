package validation

import (
	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func ParseRefreshJTI(refreshJTI string) (uuid.UUID, error) {
	jti, err := uuid.Parse(refreshJTI)
	if err != nil {
		return uuid.Nil, fun.Errf("invalid refresh jti: %s", err).BadRequest()
	}
	return jti, nil
}

func RequireRefreshJTI(refreshJTI *string) (uuid.UUID, error) {
	if refreshJTI == nil {
		return uuid.Nil, fun.Errf("%s field is nil", "refreshJTI").UnprocessableEntity()
	}
	return ParseRefreshJTI(*refreshJTI)
}
