package models

import (
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
}

type RefreshClaims struct {
	Sub RefreshSubJWT `json:"sub"`
	jwt.RegisteredClaims
}

type UserTokens struct {
	AccessTokenString  string `json:"access_token"`
	RefreshTokenString string `json:"refresh_token"`
}
