package ports

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

// PurchaseRepo is the store's persistence surface: purchases + items are the
// availability ledger (D4); materialized rows (registrations /
// product_purchases / program_participations) are written in the same tx and
// flipped by the webhook receiver (split 4) and the expiry worker (split 7).
type PurchaseRepo interface {
	// CreatePurchase inserts a pending purchase and its items. Callers wrap
	// item creation + materialized-row writes in one tx (checkout, split 7).
	CreatePurchase(ctx context.Context, toCreate *models.Purchase) (*models.Purchase, error)
	CreatePurchaseItem(ctx context.Context, toCreate *models.PurchaseItem) (*models.PurchaseItem, error)

	GetByID(ctx context.Context, id uuid.UUID) (*models.Purchase, error)
	// GetByIDForOwner is the resume (split 5) read: the authenticated user
	// must be the purchase's purchaser.
	GetByIDForOwner(ctx context.Context, id, purchaserID uuid.UUID) (*models.Purchase, error)
	// GetByIntentID is the webhook receiver's correlation read (D2).
	GetByIntentID(ctx context.Context, intentID uuid.UUID) (*models.Purchase, error)
	ListByPurchaser(ctx context.Context, purchaserID uuid.UUID) ([]models.Purchase, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.PurchaseStatus, reason *string) (*models.Purchase, error)
	ListItemsByPurchase(ctx context.Context, purchaseID uuid.UUID) ([]models.PurchaseItem, error)

	// Availability returns the stock position of every purchasable item in
	// an edition (available = base - reserved; nil base = unlimited). Used
	// by the checkout tx (split 7), the SSE relay (split 6), and the late-
	// approval path (split 4).
	Availability(ctx context.Context, editionID uuid.UUID) ([]models.ItemAvailability, error)
}
