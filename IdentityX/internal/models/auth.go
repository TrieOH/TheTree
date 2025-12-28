package models

import (
	"GoAuth/internal/apierr"
	"context"
	"encoding/json"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AccessSubJWT struct {
	ID        uuid.UUID        `json:"id"`
	Email     string           `json:"email"`
	ProjectID *uuid.UUID       `json:"projectID"`
	UserType  string           `json:"user_type"`
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

type RefreshData struct {
	RefreshCookie *http.Cookie
	Agent         string
	IP            string
}

type ctxKey string

const (
	AccessClaimsKey  ctxKey = "accessClaims"
	RefreshClaimsKey ctxKey = "refreshClaims"
)

func GetAccessClaims(ctx context.Context) (*AccessClaims, error) {
	val := ctx.Value(AccessClaimsKey)
	if val == nil {
		return nil, apierr.ErrUnauthorized.WithMsg("access token is missing or invalid").WithID(apierr.AuthMissingAccessClaims)
	}

	claims, ok := val.(*AccessClaims)
	if !ok {
		return nil, apierr.ErrInternal.WithMsg("invalid access claims type in context").WithID(apierr.AuthInvalidAccessClaims)
	}

	return claims, nil
}

func GetRefreshClaims(ctx context.Context) (*RefreshClaims, error) {
	val := ctx.Value(RefreshClaimsKey)
	if val == nil {
		return nil, apierr.ErrUnauthorized.WithMsg("refresh token is missing or invalid").WithID(apierr.AuthMissingRefreshClaims)
	}

	claims, ok := val.(*RefreshClaims)
	if !ok {
		return nil, apierr.ErrInternal.WithMsg("invalid refresh claims type in context").WithID(apierr.AuthInvalidRefreshClaims)
	}

	return claims, nil
}
