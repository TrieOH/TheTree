package utils

import (
	"GoAuth/internal/apierr"
	"GoAuth/internal/domain/auth"
	"crypto/ed25519"

	"github.com/golang-jwt/jwt/v5"
)

func ParseAccessToken(tokenStr string, secret ed25519.PublicKey) (*auth.AccessClaims, error) {
	claims := &auth.AccessClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return nil, apierr.FromJWTError(err, "access token")
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
		return nil, apierr.FromJWTError(err, "refresh token")
	}

	if token == nil || !token.Valid {
		return nil, apierr.ErrUnauthorized.WithMsg("invalid refresh token").WithID(apierr.TokenInvalid)
	}

	return claims, nil
}
