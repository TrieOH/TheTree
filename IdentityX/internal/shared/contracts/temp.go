package contracts

import (
	"time"

	"github.com/google/uuid"
)

type NewAccessTokenInput struct {
	KID                  string
	User                 User
	IP, Agent, AccessJTI string
	SessionID            uuid.UUID
	FamilyID             uuid.UUID
	ExpiresAt            time.Time
}

type NewRefreshTokenInput struct {
	KID                   string
	AccessJTI, RefreshJTI uuid.UUID
	ExpiresAt             time.Time
	FamilyID              uuid.UUID
}

type NewProjectAccessTokenInput struct {
	KID                  string
	User                 User
	IP, Agent, AccessJTI string
	SessionID            uuid.UUID
	FamilyID             uuid.UUID
	ExpiresAt            time.Time
}

type NewVerificationTokenInput struct {
	KID       string
	Subject   uuid.UUID
	ExpiresAt time.Time
}

type NewResetPasswordInput struct {
	KID       string
	Subject   uuid.UUID
	ExpiresAt time.Time
	ProjectID *uuid.UUID
}
