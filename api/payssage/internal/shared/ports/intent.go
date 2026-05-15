package ports

import (
	"context"
	"payssage/internal/shared/contracts"

	"github.com/google/uuid"
)

type IntentRepository interface {
	Create(ctx context.Context, toCreate contracts.Intent) (*contracts.Intent, error)
	GetByID(ctx context.Context, id uuid.UUID) (*contracts.Intent, error)
	List(ctx context.Context) ([]contracts.Intent, error)
	ListIntentsByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]contracts.Intent, error)
	Cancel(ctx context.Context, id uuid.UUID) (*contracts.Intent, error)
	Confirm(ctx context.Context, id uuid.UUID) (*contracts.Intent, error)
	Fail(ctx context.Context, id uuid.UUID) (*contracts.Intent, error)
	UpdateProviderData(ctx context.Context, intent contracts.Intent) (*contracts.Intent, error)
	GetByMPOrderID(ctx context.Context, orderID string) (*contracts.Intent, error)
	GetByMPTransactionID(ctx context.Context, transactionID string) (*contracts.Intent, error)
}
