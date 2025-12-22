package models

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AccessSubJWT struct {
	ID        uuid.UUID        `json:"id"`
	Email     string           `json:"email"`
	ProjectId *uuid.UUID       `json:"projectId"`
	Metadata  *json.RawMessage `json:"metadata"`
	SessionID uuid.UUID        `json:"session_id"`
	UserAgent string           `json:"user_agent"`
	UserIP    string           `json:"user_ip"`
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

type ctxKey string

const (
	AccessClaimsKey  ctxKey = "accessClaims"
	RefreshClaimsKey ctxKey = "refreshClaims"
)

func GetAccessClaims(ctx context.Context) (*AccessClaims, error) {
	val := ctx.Value(AccessClaimsKey)
	if val == nil {
		return nil, fmt.Errorf("access claims not found in context")
	}

	claims, ok := val.(*AccessClaims)
	if !ok {
		return nil, fmt.Errorf("invalid type for access claims in context")
	}

	return claims, nil
}

func GetRefreshClaims(ctx context.Context) (*RefreshClaims, error) {
	val := ctx.Value(RefreshClaimsKey)
	if val == nil {
		return nil, fmt.Errorf("refresh claims not found in context")
	}

	claims, ok := val.(*RefreshClaims)
	if !ok {
		return nil, fmt.Errorf("invalid type for refresh claims in context")
	}

	return claims, nil
}
