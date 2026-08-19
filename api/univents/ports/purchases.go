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
	// ListByEdition is the organizer read (refund plan B3): every purchase of
	// an edition, newest first — for the owner/admin orders page.
	ListByEdition(ctx context.Context, editionID uuid.UUID) ([]models.Purchase, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.PurchaseStatus, reason *string) (*models.Purchase, error)
	// UpdateStatusIf performs a guarded status transition (WHERE status =
	// from) and returns (nil, nil) when the guard misses. The webhook
	// receiver (split 4) uses it so duplicate deliveries are idempotent
	// no-ops: approve pending→approved, late-approve expired→approved,
	// failure pending→cancelled.
	UpdateStatusIf(ctx context.Context, id uuid.UUID, from, to models.PurchaseStatus, reason *string) (*models.Purchase, error)
	// UpdateRiverJobID links the expiry river job to the purchase (checkout,
	// split 7) so the webhook receiver can cancel it on approve.
	UpdateRiverJobID(ctx context.Context, id uuid.UUID, riverJobID int64) (*models.Purchase, error)
	// AttachIntent stores the Payssage intent on the purchase after the
	// intent was created (checkout, split 7): seller, intent id (the D2
	// correlation key), and the pix QR.
	AttachIntent(ctx context.Context, id uuid.UUID, sellerID, intentID uuid.UUID, qrCode, qrCodeBase64 *string) (*models.Purchase, error)
	ListItemsByPurchase(ctx context.Context, purchaseID uuid.UUID) ([]models.PurchaseItem, error)

	// Availability returns the stock position of every purchasable item in
	// an edition (available = base - reserved; nil base = unlimited). Used
	// by the checkout tx (split 7), the SSE relay (split 6), and the late-
	// approval path (split 4).
	Availability(ctx context.Context, editionID uuid.UUID) ([]models.ItemAvailability, error)
}
