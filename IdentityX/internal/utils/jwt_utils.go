package utils

import (
	"GoAuth/internal/models"
	"crypto/ed25519"
	"errors"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

func ParseAccessToken(tokenStr string, secret ed25519.PublicKey) (*models.AccessClaims, *resp.Response) {
	claims := &models.AccessClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return nil, handleJWTError(err, "access token")
	}

	if token == nil || !token.Valid {
		return nil, resp.Unauthorized("invalid access token")
	}

	return claims, nil
}

func ParseAccessTokenUserIDUnsafe(tokenStr string, secret ed25519.PublicKey) *string {
	if tokenStr == "" {
		return nil
	}

	claims := &models.AccessClaims{}

	_, _, err := new(jwt.Parser).ParseUnverified(tokenStr, claims)
	if err != nil || claims.Sub.ID == uuid.Nil {
		return nil
	}

	_, err = jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) ||
			errors.Is(err, jwt.ErrTokenNotValidYet) ||
			errors.Is(err, jwt.ErrTokenUsedBeforeIssued) {
			// These are fine - still return ID
		} else {
			// Signature invalid, malformed, etc -> token is not mine
			return nil
		}
	}

	// Success - return UserID
	idStr := claims.Sub.ID.String()
	return &idStr
}

func ParseRefreshToken(tokenStr string, secret ed25519.PublicKey) (*models.RefreshClaims, *resp.Response) {
	claims := &models.RefreshClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return nil, handleJWTError(err, "refresh token")
	}

	if token == nil || !token.Valid {
		return nil, resp.Unauthorized("invalid refresh token")
	}

	return claims, nil
}
