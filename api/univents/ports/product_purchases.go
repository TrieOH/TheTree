package ports

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

// ProductPurchaseRepo is the product_purchases write side: created pending
// at checkout (split 7) as a purchase_items materialization, flipped by the
// webhook receiver (split 4) and the expiry worker (split 7).
type ProductPurchaseRepo interface {
	Create(ctx context.Context, toCreate *models.ProductPurchase) (*models.ProductPurchase, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.ProductPurchaseStatus, reason *string) (*models.ProductPurchase, error)
}
