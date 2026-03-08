package dto

import (
	"time"

	"github.com/google/uuid"
)

type BeginOAuthResponse struct {
	RedirectURL string `json:"redirect_url"`
}

type SetMarketplaceConfigRequest struct {
	CredentialID uuid.UUID `json:"credential_id" validate:"required"`
	FeeBps       int       `json:"fee_bps" validate:"min=0,max=10000"`
}

type MarketplaceConfigResponse struct {
	ID           uuid.UUID `json:"id"`
	WorkspaceID  uuid.UUID `json:"workspace_id"`
	CredentialID uuid.UUID `json:"credential_id"`
	FeeBps       int       `json:"fee_bps"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type SetupProviderRequest struct {
	IsMarketplace    bool   `json:"is_marketplace"`
	FeeBps           int    `json:"fee_bps" validate:"min=0,max=10000"`
	FinalRedirectURL string `json:"final_redirect_url" validate:"required,url"`
}

type ConnectSellerRequest struct {
	FinalRedirectURL string `json:"final_redirect_url" validate:"required,url"`
}
