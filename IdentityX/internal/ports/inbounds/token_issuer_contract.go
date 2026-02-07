package inbounds

import (
	"GoAuth/internal/domain/project_users"
	"GoAuth/internal/domain/user"
	"time"

	"github.com/google/uuid"
)

type NewAccessTokenInput struct {
	KID                  string
	User                 user.User
	IP, Agent, AccessJTI string
	SessionID            uuid.UUID
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
	User                 project_users.ProjectUser
	IP, Agent, AccessJTI string
	SessionID            uuid.UUID
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
