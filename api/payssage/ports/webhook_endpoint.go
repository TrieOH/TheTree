package ports

import (
	"context"
	"payssage/models"

	"github.com/google/uuid"
)

type WebhookEndpointRepo interface {
	Create(ctx context.Context, toCreate models.WebhookEndpoint) (*models.WebhookEndpoint, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.WebhookEndpoint, error)
	ListByWallet(ctx context.Context, walletID uuid.UUID) ([]models.WebhookEndpoint, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
