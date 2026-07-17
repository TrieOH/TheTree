package models

import (
	"github.com/google/uuid"

	"time"
)

type OAuthFlow string

const (
	OAuthFlowCollector OAuthFlow = "collector"
	OAuthFlowSeller    OAuthFlow = "seller"
)

func (f OAuthFlow) IsValid() bool {
	switch f {
	case OAuthFlowCollector, OAuthFlowSeller:
		return true
	default:
		return false
	}
}

func (f OAuthFlow) String() string {
	return string(f)
}

type OAuthState struct {
	State            string     `json:"state"`
	WalletID         *uuid.UUID `json:"wallet_id"`
	Provider         string     `json:"provider"`
	Flow             OAuthFlow  `json:"flow"`
	FinalRedirectUrl string     `json:"final_redirect_url"`
	CreatedAt        time.Time  `json:"created_at"`
	ExpiresAt        time.Time  `json:"expires_at"`
}

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
	PublicKey      string `json:"public_key,omitempty"`
}

type ConnectRequest struct {
	Flow                OAuthFlow  `json:"flow"                  validate:"required"`
	ProviderRedirectURL string     `json:"provider_redirect_url" validate:"required,url"`
	FinalRedirectURL    string     `json:"final_redirect_url"    validate:"required,url"`
	WalletID            *uuid.UUID `json:"wallet_id"`
}

func (r *ConnectRequest) ToInput(p string) ConnectInput {
	return ConnectInput{
		Provider:            p,
		Flow:                r.Flow,
		ProviderRedirectURL: r.ProviderRedirectURL,
		FinalRedirectURL:    r.FinalRedirectURL,
		WalletID:            r.WalletID,
	}
}

type ConnectInput struct {
	Provider            string
	Flow                OAuthFlow
	ProviderRedirectURL string
	FinalRedirectURL    string
	WalletID            *uuid.UUID
}
