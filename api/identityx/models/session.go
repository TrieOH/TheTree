package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	SessionID uuid.UUID  `json:"session_id"`
	ProjectID *uuid.UUID `json:"project_id"`
	UserID    uuid.UUID  `json:"user_id"`
	UserType  UserType   `json:"user_type"`
	FamilyID  uuid.UUID  `json:"family_id"`
	TokenID   uuid.UUID  `json:"token_id"`
	IssuedAt  time.Time  `json:"issued_at"`
	UserAgent string     `json:"user_agent"`
	UserIP    string     `json:"user_ip"`
	RevokedAt *time.Time `json:"revoked_at"`
	ExpiresAt time.Time  `json:"expires_at"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type Filter struct {
	UserID    uuid.UUID  `json:"user_id"`
	ExcludeID *uuid.UUID `json:"exclude_id"`
	UserType  UserType   `json:"user_type"`
}
