package models

import (
	"encoding/json"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AccessSub struct {
	ID           uuid.UUID        `json:"id"`
	ProjectID    *uuid.UUID       `json:"project_id"`
	Email        *string          `json:"email"`
	Type         ActorType        `json:"type"`
	VerifiedAt   *time.Time       `json:"verified_at"`
	Capabilities json.RawMessage  `json:"capabilities"`
	Metadata     *json.RawMessage `json:"metadata"`
}

type AccessClaims struct {
	jwt.RegisteredClaims

	Sub AccessSub `json:"subject"`
}

type RefreshSub struct {
	ID        uuid.UUID  `json:"id"`
	ProjectID *uuid.UUID `json:"project_id"`
	AccessJTI uuid.UUID  `json:"access_jti"`
}

type RefreshClaims struct {
	jwt.RegisteredClaims

	Sub RefreshSub `json:"subject"`
}

type UserTokensOutput struct {
	AccessToken      string    `json:"access_token"`
	RefreshToken     string    `json:"refresh_token"`
	AccessExpiresAt  time.Time `json:"access_expires_at"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at"`
	Domain           string    `json:"domain"`
}

type UserTokensResponse struct {
	AccessTokenString  string    `json:"access_token_string"`
	RefreshTokenString string    `json:"refresh_token_string"`
	AccessExpiresAt    time.Time `json:"access_expires_at"`
	RefreshExpiresAt   time.Time `json:"refresh_expires_at"`
	Domain             string    `json:"domain"`
}

func (r *UserTokensOutput) ToResponse() UserTokensResponse {
	return UserTokensResponse{
		AccessTokenString:  r.AccessToken,
		RefreshTokenString: r.RefreshToken,
		AccessExpiresAt:    r.AccessExpiresAt,
		RefreshExpiresAt:   r.RefreshExpiresAt,
		Domain:             r.Domain,
	}
}
