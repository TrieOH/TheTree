package utils

import (
	"errors"

	resp "github.com/MintzyG/GoResponse/response"
	"GoAuth/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

func handleJWTError(err error, tokenType string) *resp.Response {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, jwt.ErrTokenExpired):
		return resp.Unauthorized(tokenType + " expired").AddTrace(err)
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		return resp.Unauthorized("invalid " + tokenType + " signature").AddTrace(err)
	case errors.Is(err, jwt.ErrTokenMalformed):
		return resp.Unauthorized("malformed " + tokenType).AddTrace(err)
	case errors.Is(err, jwt.ErrTokenInvalidClaims):
		return resp.Unauthorized("invalid " + tokenType + " claims").AddTrace(err)
	case errors.Is(err, jwt.ErrTokenNotValidYet):
		return resp.Unauthorized(tokenType + " not valid yet").AddTrace(err)
	case errors.Is(err, jwt.ErrTokenUsedBeforeIssued):
		return resp.Unauthorized(tokenType + " used before issued").AddTrace(err)
	case errors.Is(err, jwt.ErrTokenInvalidIssuer):
		return resp.Unauthorized(tokenType + " has invalid issuer").AddTrace(err)
	case errors.Is(err, jwt.ErrTokenInvalidSubject):
		return resp.Unauthorized(tokenType + " has invalid subject").AddTrace(err)
	case errors.Is(err, jwt.ErrTokenInvalidAudience):
		return resp.Unauthorized(tokenType + " has invalid audience").AddTrace(err)
	case errors.Is(err, jwt.ErrTokenInvalidId):
		return resp.Unauthorized(tokenType + " has invalid id").AddTrace(err)
	}

	return resp.Unauthorized("invalid " + tokenType).AddTrace(err)
}

func ParseAccessToken(tokenStr, secret string) (*models.AccessClaims, *resp.Response) {
	claims := &models.AccessClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, handleJWTError(err, "access token")
	}

	if token == nil || !token.Valid {
		return nil, resp.Unauthorized("invalid access token")
	}

	return claims, nil
}

func ParseRefreshToken(tokenStr, secret string) (*models.RefreshClaims, *resp.Response) {
	claims := &models.RefreshClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, handleJWTError(err, "refresh token")
	}

	if token == nil || !token.Valid {
		return nil, resp.Unauthorized("invalid refresh token")
	}

	return claims, nil
}
