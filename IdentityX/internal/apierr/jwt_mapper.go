package apierr

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

func FromJWTError(err error, tokenType string) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, jwt.ErrTokenExpired):
		return ErrUnauthorized.WithMsg(tokenType + " expired").WithID(TokenExpired)
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		return ErrUnauthorized.WithMsg("invalid " + tokenType + " signature").WithID(TokenSignatureInvalid)
	case errors.Is(err, jwt.ErrTokenMalformed):
		return ErrUnauthorized.WithMsg("malformed " + tokenType + " token").WithID(TokenMalformed)
	case errors.Is(err, jwt.ErrTokenInvalidClaims):
		return ErrUnauthorized.WithMsg("invalid " + tokenType + " claims").WithID(TokenInvalidAccessClaims)
	case errors.Is(err, jwt.ErrTokenNotValidYet):
		return ErrUnauthorized.WithMsg(tokenType + " not yet valid").WithID(TokenNotYetValid)
	case errors.Is(err, jwt.ErrTokenUsedBeforeIssued):
		return ErrUnauthorized.WithMsg(tokenType + " used before issued").WithID(TokenUsedBeforeIssued)
	case errors.Is(err, jwt.ErrTokenInvalidIssuer):
		return ErrUnauthorized.WithMsg(tokenType + " has invalid issuer").WithID(TokenInvalidIssuer)
	case errors.Is(err, jwt.ErrTokenInvalidSubject):
		return ErrUnauthorized.WithMsg(tokenType + " has invalid subject").WithID(TokenInvalidSubject)
	case errors.Is(err, jwt.ErrTokenInvalidAudience):
		return ErrUnauthorized.WithMsg(tokenType + " has invalid audience").WithID(TokenInvalidAudience)
	case errors.Is(err, jwt.ErrTokenInvalidId):
		if tokenType == "refresh token" {
			return ErrUnauthorized.WithMsg(tokenType + " has invalid id").WithID(TokenRefreshInvalidID)
		}
		return ErrUnauthorized.WithMsg(tokenType + " has invalid id").WithID(TokenAccessInvalidID)
	}

	return ErrUnauthorized.WithMsg("invalid " + tokenType + " token").WithID(TokenInvalid).WithCause(err)
}
