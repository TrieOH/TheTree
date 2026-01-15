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
		return ErrUnauthorized.WithMsg(tokenType + " token expired").WithID(TokenExpired)
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		return ErrUnauthorized.WithMsg("invalid " + tokenType + " token signature").WithID(TokenSignatureInvalid)
	case errors.Is(err, jwt.ErrTokenMalformed):
		return ErrUnauthorized.WithMsg("malformed " + tokenType + " token").WithID(TokenMalformed)
	case errors.Is(err, jwt.ErrTokenInvalidClaims):
		return ErrUnauthorized.WithMsg("invalid " + tokenType + " token claims").WithID(TokenInvalidAccessClaims)
	case errors.Is(err, jwt.ErrTokenNotValidYet):
		return ErrUnauthorized.WithMsg(tokenType + " token not yet valid").WithID(TokenNotYetValid)
	case errors.Is(err, jwt.ErrTokenUsedBeforeIssued):
		return ErrUnauthorized.WithMsg(tokenType + " token used before issued").WithID(TokenUsedBeforeIssued)
	case errors.Is(err, jwt.ErrTokenInvalidIssuer):
		return ErrUnauthorized.WithMsg(tokenType + " token has invalid issuer").WithID(TokenInvalidIssuer)
	case errors.Is(err, jwt.ErrTokenInvalidSubject):
		return ErrUnauthorized.WithMsg(tokenType + " token has invalid subject").WithID(TokenInvalidSubject)
	case errors.Is(err, jwt.ErrTokenInvalidAudience):
		return ErrUnauthorized.WithMsg(tokenType + " token has invalid audience").WithID(TokenInvalidAudience)
	case errors.Is(err, jwt.ErrTokenInvalidId):
		if tokenType == "refresh" {
			return ErrUnauthorized.WithMsg(tokenType + " token has invalid id").WithID(TokenRefreshInvalidID)
		}
		return ErrUnauthorized.WithMsg(tokenType + " token has invalid id").WithID(TokenAccessInvalidID)
	case errors.Is(err, jwt.ErrTokenUnverifiable):
		return ErrUnauthorized.WithMsg("unverifiable " + tokenType + " token").WithID(TokenUnverifiable)
	}

	return ErrUnauthorized.WithMsg("invalid " + tokenType + " token").WithID(TokenInvalid).WithCause(err)
}
