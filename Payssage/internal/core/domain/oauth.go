package domain

import (
	"github.com/google/uuid"

	"time"
)

type OAuthState struct {
	State            string
	WorkspaceID      uuid.UUID
	Provider         string
	Flow             string // "setup" or "connect"
	IsMarketplace    bool
	FeeBps           int
	FinalRedirectURL string
	CreatedAt        time.Time
	ExpiresAt        time.Time
}

const (
	OAuthFlowSetup   = "setup"
	OAuthFlowConnect = "connect"
)

type ProviderCredential struct {
	ID          uuid.UUID
	WorkspaceID uuid.UUID
	Provider    string
	Credentials ProviderCredentialData
	CreatedAt   time.Time
	RevokedAt   *time.Time
}

type ProviderCredentialData struct {
	AccessToken    string `json:"access_token"`
	RefreshToken   string `json:"refresh_token,omitempty"`
	ProviderUserID int    `json:"provider_user_id,omitempty"` // MP seller ID
	Nickname       string `json:"nickname,omitempty"`
}
