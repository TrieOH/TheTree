package contracts

import (
	"time"

	"github.com/google/uuid"
)

type IdentityType string

const (
	ClientIdentity  IdentityType = "client"
	ProjectIdentity IdentityType = "project"
)

type Session struct {
	SessionID  uuid.UUID  `json:"session_id"`
	ProjectID  *uuid.UUID `json:"project_id"`
	IdentityID uuid.UUID  `json:"identity_id"`
	FamilyID   uuid.UUID  `json:"family_id"`
	TokenID    uuid.UUID  `json:"token_id"`
	IssuedAt   time.Time  `json:"issued_at"`
	UserAgent  string     `json:"user_agent"`
	UserIP     string     `json:"user_ip"`
	RevokedAt  *time.Time `json:"revoked_at"`
	ExpiresAt  time.Time  `json:"expires_at"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	UserType   string     `json:"user_type"`
}

type Identity struct {
	ID           uuid.UUID    `json:"id"`
	IdentityType IdentityType `json:"identity_type"`
	EntityID     uuid.UUID    `json:"entity_id"`
	CreatedAt    time.Time    `json:"created_at"`
}

type Filter struct {
	EntityID     uuid.UUID    `json:"entity_id"`
	ExcludeID    *uuid.UUID   `json:"exclude_id"`
	IdentityType IdentityType `json:"identity_type"`
}
