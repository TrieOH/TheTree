package ports

import (
	"context"
	"payssage/models"

	"github.com/google/uuid"
)

type IntentRepo interface {
	Create(ctx context.Context, toCreate models.Intent) (*models.Intent, error)
	Update(ctx context.Context, toUpdate models.Intent) (*models.Intent, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Intent, error)
	ListByWallet(ctx context.Context, walletID uuid.UUID) ([]models.Intent, error)
	ListByOwner(ctx context.Context, ownerID uuid.UUID) ([]models.Intent, error)
	ListByOrg(ctx context.Context, orgID uuid.UUID) ([]models.Intent, error)
}
