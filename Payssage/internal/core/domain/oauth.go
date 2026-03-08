package domain

import (
	"github.com/google/uuid"

	"time"
)

type OAuthState struct {
	State            string
	WorkspaceID      uuid.UUID
	Provider         string
	FinalRedirectURL string
	CreatedAt        time.Time
	ExpiresAt        time.Time
}

type ProviderCredential struct {
	ID          uuid.UUID
	WorkspaceID uuid.UUID
	Provider    string
	DisplayName string
	Credentials ProviderCredentialData
	CreatedAt   time.Time
	RevokedAt   *time.Time
}

type ProviderCredentialData struct {
	AccessToken    string `json:"access_token"`
	RefreshToken   string `json:"refresh_token,omitempty"`
	ProviderUserID string `json:"provider_user_id,omitempty"` // MP seller ID
}
