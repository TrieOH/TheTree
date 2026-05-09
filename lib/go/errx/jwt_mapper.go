package errx

import (
	"errors"

	"github.com/MintzyG/fun"
	"github.com/golang-jwt/jwt/v5"
)

func FromJWTError(err error, tokenType string) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, jwt.ErrTokenExpired):
		return fun.Errf("%s token expired", tokenType).Unauthorized()
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		return fun.Errf("invalid %s token signature", tokenType).Unauthorized()
	case errors.Is(err, jwt.ErrTokenMalformed):
		return fun.Errf("malformed %s token", tokenType).Unauthorized()
	case errors.Is(err, jwt.ErrTokenInvalidClaims):
		return fun.Errf("invalid %s token claims", tokenType).Unauthorized()
	case errors.Is(err, jwt.ErrTokenNotValidYet):
		return fun.Errf("%s token not yet valid", tokenType).Unauthorized()
	case errors.Is(err, jwt.ErrTokenUsedBeforeIssued):
		return fun.Errf("%s token used before issued", tokenType).Unauthorized()
	case errors.Is(err, jwt.ErrTokenInvalidIssuer):
		return fun.Errf("%s token has invalid issuer", tokenType).Unauthorized()
	case errors.Is(err, jwt.ErrTokenInvalidSubject):
		return fun.Errf("%s token has invalid subject", tokenType).Unauthorized()
	case errors.Is(err, jwt.ErrTokenInvalidAudience):
		return fun.Errf("%s token has invalid audience", tokenType).Unauthorized()
	case errors.Is(err, jwt.ErrTokenInvalidId):
		return fun.Errf("%s token has invalid id", tokenType).Unauthorized()
	case errors.Is(err, jwt.ErrTokenUnverifiable):
		return fun.Errf("unverifiable %s token", tokenType).Unauthorized()
	}

	return fun.Errf("invalid %s token", tokenType).Unauthorized()
}
