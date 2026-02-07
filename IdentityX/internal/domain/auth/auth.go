package auth

import (
	"encoding/json"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AccessSub struct {
	ID         uuid.UUID        `json:"id"`
	Email      string           `json:"email"`
	ProjectID  *uuid.UUID       `json:"project_id"`
	UserType   string           `json:"user_type"`
	Metadata   *json.RawMessage `json:"metadata"`
	SessionID  uuid.UUID        `json:"session_id"`
	UserAgent  string           `json:"user_agent"`
	UserIP     string           `json:"user_ip"`
	IsVerified bool             `json:"is_verified"`
	VerifiedAt *time.Time       `json:"verified_at"`
}

type AccessClaims struct {
	Sub AccessSub `json:"sub"`
	jwt.RegisteredClaims
}

type RefreshSub struct {
	AccessJTI uuid.UUID `json:"access_jti"`
	FamilyID  uuid.UUID `json:"family_id"`
}

type RefreshClaims struct {
	Sub RefreshSub `json:"sub"`
	jwt.RegisteredClaims
}

type VerificationSub struct {
	Subject uuid.UUID `json:"subject"`
}

type VerificationClaims struct {
	Sub VerificationSub `json:"sub"`
	jwt.RegisteredClaims
}

type ResetPasswordSub struct {
	Subject   uuid.UUID  `json:"subject"`
	ProjectID *uuid.UUID `json:"project_id"`
}

type ResetPasswordClaims struct {
	Sub ResetPasswordSub `json:"sub"`
	jwt.RegisteredClaims
}
