package ports

import (
	"context"

	"payssage/internal/shared/contracts"

	"github.com/google/uuid"
)

type WebhookDispatcher interface {
	Dispatch(ctx context.Context, provider, intentID, event string, eventID uuid.UUID) error
}

type WebhookEndpointRepo interface {
	Create(ctx context.Context, toCreate contracts.WebhookEndpoint) (*contracts.WebhookEndpoint, error)
	GetByID(ctx context.Context, id uuid.UUID) (*contracts.WebhookEndpoint, error)
	ListByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]contracts.WebhookEndpoint, error)
	Delete(ctx context.Context, id, workspaceID uuid.UUID) error
}

type WebhookDeliveryRepo interface {
	Create(ctx context.Context, toCreate contracts.WebhookDelivery) (*contracts.WebhookDelivery, error)
	GetByID(ctx context.Context, id uuid.UUID) (*contracts.WebhookDelivery, error)
	ListByEndpoint(ctx context.Context, endpointID uuid.UUID) ([]contracts.WebhookDelivery, error)
	MarkDelivered(ctx context.Context, id uuid.UUID) (*contracts.WebhookDelivery, error)
	MarkFailed(ctx context.Context, id uuid.UUID) (*contracts.WebhookDelivery, error)
	IncrementAttempt(ctx context.Context, id uuid.UUID) (*contracts.WebhookDelivery, error)
}

type WebhookEventRepo interface {
	Create(ctx context.Context, toCreate contracts.WebhookEventOriginal) (*contracts.WebhookEventOriginal, error)
	Enrich(ctx context.Context, id, workspaceID, intentID uuid.UUID, externalID string) (*contracts.WebhookEventOriginal, error)
	GetByID(ctx context.Context, id uuid.UUID) (*contracts.WebhookEventOriginal, error)
	ListByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]contracts.WebhookEventOriginal, error)
	ListByProvider(ctx context.Context, provider string) ([]contracts.WebhookEventOriginal, error)
}
