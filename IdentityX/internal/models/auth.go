package models

import (
	"net/http"
	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
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
	UserAgent string `json:"user_agent"`
	UserIP string `json:"user_ip"`
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

func GetAccessClaims(r *http.Request) (*AccessClaims, bool) {
	claims, ok := r.Context().Value(AccessClaimsKey).(*AccessClaims)
	return claims, ok
}

func GetRefreshClaims(r *http.Request) (*RefreshClaims, bool) {
	claims, ok := r.Context().Value(RefreshClaimsKey).(*RefreshClaims)
	return claims, ok
}
