package purchases

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"univents/internal/sqlc"
	"univents/models"
)

func (repo *Repo) CreatePurchase(ctx context.Context, toCreate *models.Purchase) (*models.Purchase, error) {
	ctx, span := telemetry.StartSpan(ctx, "PurchasesRepo.CreatePurchase")
	defer span.End()
	result, err := database.Queries(ctx, repo.q).CreatePurchase(ctx, sqlc.CreatePurchaseParams{
		EditionID:        toCreate.EditionID,
		PurchaserID:      toCreate.PurchaserID,
		Status:           string(toCreate.Status),
		StatusReason:     toCreate.StatusReason,
		TotalCents:       toCreate.TotalCents,
		Currency:         toCreate.Currency,
		PaymentMethod:    toCreate.PaymentMethod,
		PayssageSellerID: toCreate.PayssageSellerID,
		PayssageIntentID: toCreate.PayssageIntentID,
		QrCode:           toCreate.QRCode,
		QrCodeBase64:     toCreate.QRCodeBase64,
		ExpiresAt:        toCreate.ExpiresAt,
		RiverJobID:       toCreate.RiverJobID,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapPurchase(result)), nil
}

func (repo *Repo) CreatePurchaseItem(ctx context.Context, toCreate *models.PurchaseItem) (*models.PurchaseItem, error) {
	ctx, span := telemetry.StartSpan(ctx, "PurchasesRepo.CreatePurchaseItem")
	defer span.End()
	result, err := database.Queries(ctx, repo.q).CreatePurchaseItem(ctx, sqlc.CreatePurchaseItemParams{
		PurchaseID:        toCreate.PurchaseID,
		ItemType:          string(toCreate.ItemType),
		ItemID:            toCreate.ItemID,
		Quantity:          toCreate.Quantity,
		UnitPriceCents:    toCreate.UnitPriceCents,
		RegistrationID:    toCreate.RegistrationID,
		ProductPurchaseID: toCreate.ProductPurchaseID,
		ParticipationID:   toCreate.ParticipationID,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapPurchaseItem(result)), nil
}
