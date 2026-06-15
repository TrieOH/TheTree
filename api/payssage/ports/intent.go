package ports

import (
	"context"

	"payssage/models"

	"github.com/google/uuid"
)

type IntentRepository interface {
	Create(ctx context.Context, toCreate models.Intent) (*models.Intent, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Intent, error)
	List(ctx context.Context) ([]models.Intent, error)
	ListIntentsByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]models.Intent, error)
	Cancel(ctx context.Context, id uuid.UUID) (*models.Intent, error)
	Confirm(ctx context.Context, id uuid.UUID) (*models.Intent, error)
	Fail(ctx context.Context, id uuid.UUID) (*models.Intent, error)
	UpdateProviderData(ctx context.Context, intent models.Intent) (*models.Intent, error)
	GetByMPOrderID(ctx context.Context, orderID string) (*models.Intent, error)
	GetByMPTransactionID(ctx context.Context, transactionID string) (*models.Intent, error)
}
