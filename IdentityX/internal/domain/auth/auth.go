package auth

import (
	"encoding/json"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AccessSubJWT struct {
	ID        uuid.UUID        `json:"id"`
	Email     string           `json:"email"`
	ProjectID *uuid.UUID       `json:"project_id"`
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
