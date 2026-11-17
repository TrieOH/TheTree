package service

import (
	"context"
	"lib/telemetry"
	"log/slog"
	"univents/models"
)

// revert cancels a purchase whose intent could not be created (rejected
// card, provider down, revoked seller…): guarded pending→cancelled,
// materialized rows flipped to cancelled (stock freed), NOTIFY. Best-effort:
// the original payment error is what the buyer sees; failures here are
// logged (the expiry worker is the backstop).
func (o *Operations) revert(ctx context.Context, purchase *models.Purchase, items []models.PurchaseItem) {
	err := o.Tx.WithinTx(ctx, func(ctx context.Context) error {
		updated, err := o.Purchases.UpdateStatusIf(ctx, purchase.ID,
			models.PurchaseStatusPending, models.PurchaseStatusCancelled, nil)
		if err != nil {
			return err
		}
		if updated == nil {
			return nil // already flipped (e.g. a webhook won the race)
		}
		return o.flipCancelled(ctx, items)
	})
	if err != nil {
		telemetry.Log().Error("checkout: revert failed after intent creation error",
			slog.String("purchase_id", purchase.ID.String()),
			telemetry.Err(err))
		return
	}
	o.notifyStock(ctx, purchase, items)
	o.notifyPurchase(ctx, purchase, models.PurchaseStatusCancelled)
}

// flipCancelled materializes the revert (D4), inside the revert tx:
// registrations/product_purchases/participations → cancelled.
func (o *Operations) flipCancelled(ctx context.Context, items []models.PurchaseItem) error {
	for _, item := range items {
		switch item.ItemType {
		case models.PurchaseItemTypeTicket:
			if item.RegistrationID == nil {
				continue
			}
			_, err := o.Registrations.UpdateStatus(ctx, *item.RegistrationID, models.RegistrationStatusCancelled, nil)
			if err != nil {
				return err
			}
		case models.PurchaseItemTypeProduct:
			if item.ProductPurchaseID == nil {
				continue
			}
			_, err := o.ProductPurchases.UpdateProductPurchaseStatus(ctx, *item.ProductPurchaseID, models.ProductPurchaseStatusCancelled, nil)
			if err != nil {
				return err
			}
		case models.PurchaseItemTypeProgramOccurrence:
			if item.ParticipationID == nil {
				continue
			}
			_, err := o.Participations.UpdateParticipationStatus(ctx, *item.ParticipationID, models.ProgramParticipationStatusCancelled)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
