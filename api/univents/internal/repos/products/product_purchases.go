package products

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"univents/internal/sqlc"
	"univents/models"

	"github.com/google/uuid"
)

// CreateProductPurchase inserts a pending product purchase at checkout
// (split 7) — one of the purchase's materialized rows (D4).
func (repo *Repo) CreateProductPurchase(ctx context.Context, toCreate *models.ProductPurchase) (*models.ProductPurchase, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProductsRepo.CreateProductPurchase")
	defer span.End()
	result, err := database.Queries(ctx, repo.q).CreateProductPurchase(ctx, sqlc.CreateProductPurchaseParams{
		EditionID:        toCreate.EditionID,
		VariantID:        toCreate.VariantID,
		PurchaserID:      toCreate.PurchaserID,
		RegistrationID:   toCreate.RegistrationID,
		Quantity:         toCreate.Quantity,
		Status:           string(toCreate.Status),
		StatusReason:     toCreate.StatusReason,
		PayssageIntentID: toCreate.PayssageIntentID,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapProductPurchase(result)), nil
}

// UpdateProductPurchaseStatus flips a product purchase (pending→confirmed/
// cancelled/expired) — the webhook receiver (split 4) and the expiry worker
// (split 7) write side.
func (repo *Repo) UpdateProductPurchaseStatus(ctx context.Context, id uuid.UUID, status models.ProductPurchaseStatus, reason *string) (*models.ProductPurchase, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProductsRepo.UpdateProductPurchaseStatus")
	defer span.End()
	result, err := database.Queries(ctx, repo.q).UpdateProductPurchaseStatus(ctx, sqlc.UpdateProductPurchaseStatusParams{
		ID:           id,
		Status:       string(status),
		StatusReason: reason,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapProductPurchase(result)), nil
}
