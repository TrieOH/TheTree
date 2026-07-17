package ports

import (
	"context"
	"payssage/models"

	"github.com/google/uuid"
)

type SellerRepo interface {
	Create(ctx context.Context, toCreate models.Seller) (*models.Seller, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Seller, error)
	List(ctx context.Context) ([]models.Seller, error)
	ListByWallet(ctx context.Context, walletID uuid.UUID) ([]models.Seller, error)
	Revoke(ctx context.Context, id uuid.UUID) error
}
