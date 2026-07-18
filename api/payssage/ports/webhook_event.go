package ports

import (
	"context"
	"payssage/models"

	"github.com/google/uuid"
)

type WebhookEventRepo interface {
	Create(ctx context.Context, toCreate models.WebhookEvent) (*models.WebhookEvent, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.WebhookEvent, error)
	ListByWallet(ctx context.Context, walletID uuid.UUID) ([]models.WebhookEvent, error)
	ListByProvider(ctx context.Context, provider string) ([]models.WebhookEvent, error)
}
