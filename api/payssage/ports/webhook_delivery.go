package ports

import (
	"context"
	"payssage/models"
	"time"

	"github.com/google/uuid"
)

type WebhookDeliveryRepo interface {
	Create(ctx context.Context, toCreate models.WebhookDelivery) (*models.WebhookDelivery, error)
	Update(ctx context.Context, params models.UpdateDeliveryParams) (*models.WebhookDelivery, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.WebhookDelivery, error)
	ListByEndpoint(ctx context.Context, endpointID uuid.UUID) ([]models.WebhookDelivery, error)
	MarkDelivered(ctx context.Context, id uuid.UUID) (*models.WebhookDelivery, error)
	MarkFailed(ctx context.Context, id uuid.UUID) (*models.WebhookDelivery, error)
	IncrementAttempt(ctx context.Context, id uuid.UUID) (*models.WebhookDelivery, error)
	// ListUnsent returns pending deliveries never attempted and older than
	// olderThan, oldest first, at most limit of them.
	ListUnsent(ctx context.Context, olderThan time.Duration, limit int) ([]models.WebhookDelivery, error)
}
