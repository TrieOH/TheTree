package ports

import (
	"context"
	"payssage/models"

	"github.com/google/uuid"
)

type WalletRepo interface {
	Create(ctx context.Context, toCreate models.Wallet) (*models.Wallet, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Wallet, error)
	GetRole(ctx context.Context, actorID, walletID uuid.UUID) (models.OrganizationRole, error)
	List(ctx context.Context, ownerID uuid.UUID) ([]models.Wallet, error)
	ListFromOrg(ctx context.Context, orgID uuid.UUID) ([]models.Wallet, error)
	SetSandboxState(ctx context.Context, walletID uuid.UUID, state bool) error
	SetFeeBPS(ctx context.Context, walletID uuid.UUID, feeBPS int) error
	BindCollector(ctx context.Context, walletID, collectorID uuid.UUID) error
	UnbindCollector(ctx context.Context, walletID uuid.UUID) error
}
