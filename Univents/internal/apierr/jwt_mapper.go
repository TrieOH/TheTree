package apierr

import (
	"errors"

	"github.com/MintzyG/fail/v3"
	"github.com/golang-jwt/jwt/v5"
)

// FIXME transform this into a fail.Mapper

func FromJWTError(err error, tokenType string) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, jwt.ErrTokenExpired):
		return fail.New(TokenExpired).WithArgs(tokenType)
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		return fail.New(TokenSignatureInvalid).WithArgs(tokenType)
	case errors.Is(err, jwt.ErrTokenMalformed):
		return fail.New(TokenMalformed).WithArgs(tokenType)
	case errors.Is(err, jwt.ErrTokenInvalidClaims):
		return fail.New(TokenInvalidAccessClaims).WithArgs(tokenType)
	case errors.Is(err, jwt.ErrTokenNotValidYet):
		return fail.New(TokenNotYetValid).WithArgs(tokenType)
	case errors.Is(err, jwt.ErrTokenUsedBeforeIssued):
		return fail.New(TokenUsedBeforeIssued).WithArgs(tokenType)
	case errors.Is(err, jwt.ErrTokenInvalidIssuer):
		return fail.New(TokenInvalidIssuer).WithArgs(tokenType)
	case errors.Is(err, jwt.ErrTokenInvalidSubject):
		return fail.New(TokenInvalidSubject).WithArgs(tokenType)
	case errors.Is(err, jwt.ErrTokenInvalidAudience):
		return fail.New(TokenInvalidAudience).WithArgs(tokenType)
	case errors.Is(err, jwt.ErrTokenInvalidId):
		if tokenType == "refresh" {
			return fail.New(TokenRefreshInvalidID).WithArgs(tokenType).Trace("couldn't parse refresh token JTI")
		}
		return fail.New(TokenAccessInvalidID).WithArgs(tokenType).Trace("couldn't parse access token JTI")
	case errors.Is(err, jwt.ErrTokenUnverifiable):
		return fail.New(TokenUnverifiable).WithArgs(tokenType)
	}

	return fail.New(TokenInvalid).WithArgs(tokenType).Trace(err.Error())
}
