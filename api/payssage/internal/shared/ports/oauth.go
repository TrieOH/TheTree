package ports

import (
	"context"

	"payssage/internal/shared/contracts"

	"github.com/google/uuid"
)

type OAuthStateRepo interface {
	Create(ctx context.Context, state contracts.OAuthState) (*contracts.OAuthState, error)
	Get(ctx context.Context, state string) (*contracts.OAuthState, error)
	Delete(ctx context.Context, state string) error
}

type ProviderCredentialRepo interface {
	Create(ctx context.Context, cred contracts.ProviderCredential) (*contracts.ProviderCredential, error)
	GetByID(ctx context.Context, id uuid.UUID) (*contracts.ProviderCredential, error)
	ListByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]contracts.ProviderCredential, error)
	Revoke(ctx context.Context, id uuid.UUID, workspaceID uuid.UUID) (*contracts.ProviderCredential, error)
	GetByWorkspaceAndProvider(ctx context.Context, workspaceID uuid.UUID, provider string) (*contracts.ProviderCredential, error)
	GetSellerCredentialByProvider(ctx context.Context, workspaceID uuid.UUID, provider string) (*contracts.ProviderCredential, error)
}

type MarketplaceConfigRepo interface {
	Create(ctx context.Context, config contracts.MarketplaceConfig) (*contracts.MarketplaceConfig, error)
	List(ctx context.Context, workspaceID uuid.UUID) ([]contracts.MarketplaceConfig, error)
	Get(ctx context.Context, workspaceID, credentialID uuid.UUID) (*contracts.MarketplaceConfig, error)
	GetByProvider(ctx context.Context, workspaceID uuid.UUID, provider string) (*contracts.MarketplaceConfig, error)
	Update(ctx context.Context, config contracts.MarketplaceConfig) (*contracts.MarketplaceConfig, error)
	Delete(ctx context.Context, workspaceID, credentialID uuid.UUID) error
	DeleteAll(ctx context.Context, workspaceID uuid.UUID) error
}

type OAuthProvider interface {
	BuildAuthURL(state, redirectURI string) string
	ExchangeCode(ctx context.Context, code, redirectURI string) (contracts.ProviderCredentialData, error)
}
