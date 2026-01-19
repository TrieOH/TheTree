package inbounds

import (
	"GoAuth/internal/domain/project_users"
	"GoAuth/internal/domain/user"
	"crypto/ed25519"
	"time"

	"github.com/google/uuid"
)

type NewAccessTokenInput struct {
	User                        user.User
	PrivateKey                  ed25519.PrivateKey
	IP, Agent, AccessJTI, KeyID string
	SessionID                   uuid.UUID
	ExpiresAt                   time.Time
}

type NewRefreshTokenInput struct {
	KeyID                 string
	PrivateKey            ed25519.PrivateKey
	AccessJTI, RefreshJTI uuid.UUID
	ExpiresAt             time.Time
	FamilyID              uuid.UUID
}

type NewProjectAccessTokenInput struct {
	User                        project_users.ProjectUser
	IP, Agent, AccessJTI, KeyID string
	SessionID                   uuid.UUID
	ExpiresAt                   time.Time
	PrivateKey                  ed25519.PrivateKey
}
