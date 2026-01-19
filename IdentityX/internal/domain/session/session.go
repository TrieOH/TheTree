package session

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
	SessionID  uuid.UUID
	ProjectID  *uuid.UUID
	IdentityID uuid.UUID
	FamilyID   uuid.UUID
	TokenID    uuid.UUID
	IssuedAt   time.Time
	UserAgent  string
	UserIP     string
	RevokedAt  *time.Time
	ExpiresAt  time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
	UserType   string
}

type Identity struct {
	ID           uuid.UUID
	IdentityType IdentityType
	EntityID     uuid.UUID
	CreatedAt    time.Time
}

type Filter struct {
	EntityID     uuid.UUID
	ExcludeID    *uuid.UUID
	IdentityType IdentityType
}
