package domain

import (
	"context"

	"github.com/google/uuid"
)

type WebhookDispatcher interface {
	Dispatch(ctx context.Context, provider, intentID, event string) error
}

type IntentRepository interface {
	Create(ctx context.Context, toCreate Intent) (*Intent, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Intent, error)
	List(ctx context.Context) ([]Intent, error)
	ListIntentsByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]Intent, error)
	Cancel(ctx context.Context, id uuid.UUID) (*Intent, error)
	Confirm(ctx context.Context, id uuid.UUID) (*Intent, error)
	Fail(ctx context.Context, id uuid.UUID) (*Intent, error)
	Pay(ctx context.Context, id uuid.UUID, providerPaymentID string, status IntentStatus) (*Intent, error)
	GetByProviderPaymentID(ctx context.Context, providerPaymentID string) (*Intent, error)
}

type WorkspaceRepo interface {
	Create(ctx context.Context, toCreate Workspace) (*Workspace, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Workspace, error)
	GetByName(ctx context.Context, name string, userID uuid.UUID) (*Workspace, error)
	List(ctx context.Context, userID uuid.UUID) ([]Workspace, error)
	EnableSandbox(ctx context.Context, id uuid.UUID) (*Workspace, error)
	DisableSandbox(ctx context.Context, id uuid.UUID) (*Workspace, error)
}

type ApiKeysRepo interface {
	Create(ctx context.Context, toCreate APIKey) (*APIKey, error)
	GetByPrefix(ctx context.Context, prefix string) ([]APIKey, error)
	ListByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]APIKey, error)
	Revoke(ctx context.Context, id, workspaceID uuid.UUID) (*APIKey, error)
}

type WebhookEndpointRepo interface {
	Create(ctx context.Context, toCreate WebhookEndpoint) (*WebhookEndpoint, error)
	GetByID(ctx context.Context, id uuid.UUID) (*WebhookEndpoint, error)
	ListByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]WebhookEndpoint, error)
	Delete(ctx context.Context, id, workspaceID uuid.UUID) error
}

type WebhookDeliveryRepo interface {
	Create(ctx context.Context, toCreate WebhookDelivery) (*WebhookDelivery, error)
	GetByID(ctx context.Context, id uuid.UUID) (*WebhookDelivery, error)
	ListByEndpoint(ctx context.Context, endpointID uuid.UUID) ([]WebhookDelivery, error)
	MarkDelivered(ctx context.Context, id uuid.UUID) (*WebhookDelivery, error)
	MarkFailed(ctx context.Context, id uuid.UUID) (*WebhookDelivery, error)
	IncrementAttempt(ctx context.Context, id uuid.UUID) (*WebhookDelivery, error)
}

type OAuthStateRepo interface {
	Create(ctx context.Context, state OAuthState) (*OAuthState, error)
	Get(ctx context.Context, state string) (*OAuthState, error)
	Delete(ctx context.Context, state string) error
}

type ProviderCredentialRepo interface {
	Create(ctx context.Context, cred ProviderCredential) (*ProviderCredential, error)
	GetByID(ctx context.Context, id uuid.UUID) (*ProviderCredential, error)
	ListByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]ProviderCredential, error)
	Revoke(ctx context.Context, id uuid.UUID, workspaceID uuid.UUID) (*ProviderCredential, error)
	GetByWorkspaceAndProvider(ctx context.Context, workspaceID uuid.UUID, provider string) (*ProviderCredential, error)
}

type MarketplaceConfigRepo interface {
	Create(ctx context.Context, config MarketplaceConfig) (*MarketplaceConfig, error)
	Get(ctx context.Context, workspaceID uuid.UUID) (*MarketplaceConfig, error)
	Update(ctx context.Context, config MarketplaceConfig) (*MarketplaceConfig, error)
	Delete(ctx context.Context, workspaceID uuid.UUID) error
}
