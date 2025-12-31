package utils

import (
	"GoAuth/internal/apierr"
	"GoAuth/internal/domain/auth"
	"crypto/ed25519"
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

func handleJWTError(err error, tokenType string) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, jwt.ErrTokenExpired):
		return apierr.ErrUnauthorized.WithMsg(tokenType + " expired").WithID(apierr.TokenExpired)
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		return apierr.ErrUnauthorized.WithMsg("invalid " + tokenType + " signature").WithID(apierr.TokenSignatureInvalid)
	case errors.Is(err, jwt.ErrTokenMalformed):
		return apierr.ErrUnauthorized.WithMsg("malformed " + tokenType + " token").WithID(apierr.TokenMalformed)
	case errors.Is(err, jwt.ErrTokenInvalidClaims):
		return apierr.ErrUnauthorized.WithMsg("invalid " + tokenType + " claims").WithID(apierr.TokenInvalidClaims)
	case errors.Is(err, jwt.ErrTokenNotValidYet):
		return apierr.ErrUnauthorized.WithMsg(tokenType + " not yet valid").WithID(apierr.TokenNotYetValid)
	case errors.Is(err, jwt.ErrTokenUsedBeforeIssued):
		return apierr.ErrUnauthorized.WithMsg(tokenType + " used before issued").WithID(apierr.TokenUsedBeforeIssued)
	case errors.Is(err, jwt.ErrTokenInvalidIssuer):
		return apierr.ErrUnauthorized.WithMsg(tokenType + " has invalid issuer").WithID(apierr.TokenInvalidIssuer)
	case errors.Is(err, jwt.ErrTokenInvalidSubject):
		return apierr.ErrUnauthorized.WithMsg(tokenType + " has invalid subject").WithID(apierr.TokenInvalidSubject)
	case errors.Is(err, jwt.ErrTokenInvalidAudience):
		return apierr.ErrUnauthorized.WithMsg(tokenType + " has invalid audience").WithID(apierr.TokenInvalidAudience)
	case errors.Is(err, jwt.ErrTokenInvalidId):
		return apierr.ErrUnauthorized.WithMsg(tokenType + " has invalid id").WithID(apierr.TokenInvalidID)
	}

	return apierr.ErrUnauthorized.WithMsg("invalid " + tokenType + " token").WithID(apierr.TokenInvalid).WithCause(err)
}

func ParseAccessToken(tokenStr string, secret ed25519.PublicKey) (*auth.AccessClaims, error) {
	claims := &auth.AccessClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return nil, handleJWTError(err, "access token")
	}

	if token == nil || !token.Valid {
		return nil, apierr.ErrUnauthorized.WithMsg("invalid access token").WithID(apierr.TokenInvalid)
	}

	return claims, nil
}

func ParseRefreshToken(tokenStr string, secret ed25519.PublicKey) (*auth.RefreshClaims, error) {
	claims := &auth.RefreshClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return nil, handleJWTError(err, "refresh token")
	}

	if token == nil || !token.Valid {
		return nil, apierr.ErrUnauthorized.WithMsg("invalid refresh token").WithID(apierr.TokenInvalid)
	}

	return claims, nil
}
