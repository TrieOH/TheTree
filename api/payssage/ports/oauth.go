package ports

import (
	"context"

	"payssage/models"

	"github.com/google/uuid"
)

type OAuthStateRepo interface {
	Create(ctx context.Context, state models.OAuthState) (*models.OAuthState, error)
	Get(ctx context.Context, state string) (*models.OAuthState, error)
	Delete(ctx context.Context, state string) error
}

type ProviderCredentialRepo interface {
	Create(ctx context.Context, cred models.ProviderCredential) (*models.ProviderCredential, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.ProviderCredential, error)
	ListByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]models.ProviderCredential, error)
	Revoke(ctx context.Context, id uuid.UUID, workspaceID uuid.UUID) (*models.ProviderCredential, error)
	GetByWorkspaceAndProvider(ctx context.Context, workspaceID uuid.UUID, provider string) (*models.ProviderCredential, error)
	GetSellerCredentialByProvider(ctx context.Context, workspaceID uuid.UUID, provider string) (*models.ProviderCredential, error)
}

type MarketplaceConfigRepo interface {
	Create(ctx context.Context, config models.MarketplaceConfig) (*models.MarketplaceConfig, error)
	List(ctx context.Context, workspaceID uuid.UUID) ([]models.MarketplaceConfig, error)
	Get(ctx context.Context, workspaceID, credentialID uuid.UUID) (*models.MarketplaceConfig, error)
	GetByProvider(ctx context.Context, workspaceID uuid.UUID, provider string) (*models.MarketplaceConfig, error)
	Update(ctx context.Context, config models.MarketplaceConfig) (*models.MarketplaceConfig, error)
	Delete(ctx context.Context, workspaceID, credentialID uuid.UUID) error
	DeleteAll(ctx context.Context, workspaceID uuid.UUID) error
}

type OAuthProvider interface {
	BuildAuthURL(state, redirectURI string) string
	ExchangeCode(ctx context.Context, code, redirectURI string) (models.ProviderCredentialData, error)
}
