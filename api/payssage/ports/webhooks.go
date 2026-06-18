package ports

import (
	"context"

	"payssage/models"

	"github.com/google/uuid"
)

type WebhookDispatcher interface {
	Dispatch(ctx context.Context, provider, intentID, event string, eventID uuid.UUID) error
}

type WebhookEndpointRepo interface {
	Create(ctx context.Context, toCreate models.WebhookEndpoint) (*models.WebhookEndpoint, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.WebhookEndpoint, error)
	ListByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]models.WebhookEndpoint, error)
	Delete(ctx context.Context, id, workspaceID uuid.UUID) error
}

type WebhookDeliveryRepo interface {
	Create(ctx context.Context, toCreate models.WebhookDelivery) (*models.WebhookDelivery, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.WebhookDelivery, error)
	ListByEndpoint(ctx context.Context, endpointID uuid.UUID) ([]models.WebhookDelivery, error)
	MarkDelivered(ctx context.Context, id uuid.UUID) (*models.WebhookDelivery, error)
	MarkFailed(ctx context.Context, id uuid.UUID) (*models.WebhookDelivery, error)
	IncrementAttempt(ctx context.Context, id uuid.UUID) (*models.WebhookDelivery, error)
}

type WebhookEventRepo interface {
	Create(ctx context.Context, toCreate models.WebhookEventOriginal) (*models.WebhookEventOriginal, error)
	Enrich(ctx context.Context, id, workspaceID, intentID uuid.UUID, externalID string) (*models.WebhookEventOriginal, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.WebhookEventOriginal, error)
	ListByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]models.WebhookEventOriginal, error)
	ListByProvider(ctx context.Context, provider string) ([]models.WebhookEventOriginal, error)
}
