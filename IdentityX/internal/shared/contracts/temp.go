package contracts

import (
	"time"

	"github.com/google/uuid"
)

type NewAccessTokenInput struct {
	KID       string    `json:"kid"`
	User      User      `json:"user"`
	IP        string    `json:"ip"`
	Agent     string    `json:"agent"`
	AccessJTI string    `json:"access_jti"`
	SessionID uuid.UUID `json:"session_id"`
	FamilyID  uuid.UUID `json:"family_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

type NewRefreshTokenInput struct {
	KID        string    `json:"kid"`
	AccessJTI  uuid.UUID `json:"access_jti"`
	RefreshJTI uuid.UUID `json:"refresh_jti"`
	ExpiresAt  time.Time `json:"expires_at"`
	FamilyID   uuid.UUID `json:"family_id"`
}

type NewProjectAccessTokenInput struct {
	KID       string    `json:"kid"`
	User      User      `json:"user"`
	IP        string    `json:"ip"`
	Agent     string    `json:"agent"`
	AccessJTI string    `json:"access_jti"`
	SessionID uuid.UUID `json:"session_id"`
	FamilyID  uuid.UUID `json:"family_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

type NewVerificationTokenInput struct {
	KID       string    `json:"kid"`
	Subject   uuid.UUID `json:"subject"`
	ExpiresAt time.Time `json:"expires_at"`
}

type NewResetPasswordInput struct {
	KID       string     `json:"kid"`
	Subject   uuid.UUID  `json:"subject"`
	ExpiresAt time.Time  `json:"expires_at"`
	ProjectID *uuid.UUID `json:"project_id"`
}
