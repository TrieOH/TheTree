package models

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AccessSubJWT struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}

type AccessClaims struct {
	Sub AccessSubJWT `json:"sub"`
	jwt.RegisteredClaims
}

type RefreshSubJWT struct {
	AccessJTI uuid.UUID `json:"access_jti"`
	SessionID uuid.UUID `json:"session_id"`
	UserAgent string    `json:"user_agent"`
	UserIP    string    `json:"user_ip"`
}

type RefreshClaims struct {
	Sub RefreshSubJWT `json:"sub"`
	jwt.RegisteredClaims
}

type UserTokens struct {
	AccessTokenString  string `json:"access_token"`
	RefreshTokenString string `json:"refresh_token"`
}

type ctxKey string

const (
	AccessClaimsKey  ctxKey = "accessClaims"
	RefreshClaimsKey ctxKey = "refreshClaims"
)

func GetAccessClaims(r *http.Request) (*AccessClaims, error) {
	val := r.Context().Value(AccessClaimsKey)
	if val == nil {
		return nil, fmt.Errorf("access claims not found in context")
	}

	claims, ok := val.(*AccessClaims)
	if !ok {
		return nil, fmt.Errorf("invalid type for access claims in context")
	}

	return claims, nil
}

func GetRefreshClaims(r *http.Request) (*RefreshClaims, error) {
	val := r.Context().Value(RefreshClaimsKey)
	if val == nil {
		return nil, fmt.Errorf("refresh claims not found in context")
	}

	claims, ok := val.(*RefreshClaims)
	if !ok {
		return nil, fmt.Errorf("invalid type for refresh claims in context")
	}

	return claims, nil
}
