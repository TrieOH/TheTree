package ports

import (
	"context"
	"encoding/json"
	"payssage/models"

	"github.com/google/uuid"
)

type IntentRepo interface {
	Create(ctx context.Context, toCreate models.Intent) (*models.Intent, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Intent, error)
	ListByWallet(ctx context.Context, walletID uuid.UUID) ([]models.Intent, error)
	ListByOwner(ctx context.Context, ownerID uuid.UUID) ([]models.Intent, error)
	ListByOrg(ctx context.Context, orgID uuid.UUID) ([]models.Intent, error)
	Cancel(ctx context.Context, id uuid.UUID) (*models.Intent, error)
	Confirm(ctx context.Context, id uuid.UUID) (*models.Intent, error)
	Fail(ctx context.Context, id uuid.UUID) (*models.Intent, error)
	UpdateProviderData(ctx context.Context, id uuid.UUID, providerData json.RawMessage) (*models.Intent, error)
}
