package dto

import (
	"TriePayments/internal/core/domain"
	"time"

	"github.com/google/uuid"
)

type BeginOAuthResponse struct {
	RedirectURL      string `json:"redirect_url"`
	FinalRedirectURL string `json:"final_redirect_url"`
}

type SetMarketplaceConfigRequest struct {
	CredentialID uuid.UUID `json:"credential_id" validate:"required"`
	FeeBps       int       `json:"fee_bps" validate:"min=0,max=10000"`
}

type MarketplaceConfigResponse struct {
	ID           uuid.UUID `json:"id"`
	WorkspaceID  uuid.UUID `json:"workspace_id"`
	Provider     string    `json:"provider"`
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
	ProviderRedirectURL string `json:"provider_redirect_url" validate:"required,url"`
	FinalRedirectURL    string `json:"final_redirect_url" validate:"required,url"`
}

type ProviderCredentialResponse struct {
	ID          uuid.UUID  `json:"id"`
	WorkspaceID uuid.UUID  `json:"workspace_id"`
	Provider    string     `json:"provider"`
	DisplayName string     `json:"display_name"`
	CreatedAt   time.Time  `json:"created_at"`
	RevokedAt   *time.Time `json:"revoked_at"`
}

func MapProviderCredentialResponse(c *domain.ProviderCredential) ProviderCredentialResponse {
	return ProviderCredentialResponse{
		ID:          c.ID,
		WorkspaceID: c.WorkspaceID,
		Provider:    c.Provider,
		DisplayName: c.DisplayName,
		CreatedAt:   c.CreatedAt,
		RevokedAt:   c.RevokedAt,
	}
}
