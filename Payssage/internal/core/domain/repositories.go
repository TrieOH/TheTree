package domain

import (
	"context"

	"github.com/google/uuid"
)

type IntentRepository interface {
	Create(ctx context.Context, toCreate Intent) (*Intent, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Intent, error)
	List(ctx context.Context) ([]Intent, error)
	ListIntentsByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]Intent, error)
	Cancel(ctx context.Context, id uuid.UUID) (*Intent, error)
	Confirm(ctx context.Context, id uuid.UUID) (*Intent, error)
	Fail(ctx context.Context, id uuid.UUID) (*Intent, error)
}

type WorkspaceRepo interface {
	Create(ctx context.Context, toCreate Workspace) (*Workspace, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Workspace, error)
	GetByName(ctx context.Context, name string, userID uuid.UUID) (*Workspace, error)
	List(ctx context.Context, userID uuid.UUID) ([]Workspace, error)
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
