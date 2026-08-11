package webhooks

import (
	"context"

	"lib/telemetry"
	"univents/models"

	"go.uber.org/zap"
)

// cancel handles a terminal failure webhook (failed/rejected/cancelled): in
// a tx, guarded pending→cancelled (a stale failure delivery can't cancel an
// approved purchase), materialized rows flipped to cancelled, NOTIFY (stock
// freed + purchase event). Non-pending purchases are left alone and
// acknowledged.
func (o *Operations) cancel(ctx context.Context, purchase *models.Purchase) error {
	if purchase.Status != models.PurchaseStatusPending {
		telemetry.Log().Info("webhook: failure for non-pending purchase ignored",
			zap.String("purchase_id", purchase.ID.String()),
			zap.String("status", string(purchase.Status)))
		return nil
	}
	return o.tx.WithinTx(ctx, func(ctx context.Context) error {
		updated, err := o.purchases.UpdateStatusIf(ctx, purchase.ID,
			models.PurchaseStatusPending, models.PurchaseStatusCancelled, nil)
		if err != nil {
			return err
		}
		if updated == nil {
			return nil // guard missed — already flipped by a concurrent delivery
		}
		return o.flipCancelled(ctx, updated)
	})
}

// flipCancelled materializes the cancellation (D4), inside the cancel tx:
// registrations/product_purchases/participations flipped to cancelled, then
// NOTIFY (stock freed + purchase.cancelled).
func (o *Operations) flipCancelled(ctx context.Context, purchase *models.Purchase) error {
	items, err := o.purchases.ListItemsByPurchase(ctx, purchase.ID)
	if err != nil {
		return err
	}

	for _, item := range items {
		switch item.ItemType {
		case models.PurchaseItemTypeTicket:
			if item.RegistrationID == nil {
				telemetry.Log().Warn("webhook: ticket item without registration",
					zap.String("purchase_id", purchase.ID.String()))
				continue
			}
			_, err := o.registrations.UpdateStatus(ctx, *item.RegistrationID, models.RegistrationStatusCancelled, nil)
			if err != nil {
				return err
			}
		case models.PurchaseItemTypeProduct:
			if item.ProductPurchaseID == nil {
				telemetry.Log().Warn("webhook: product item without product_purchase",
					zap.String("purchase_id", purchase.ID.String()))
				continue
			}
			_, err := o.productPurchases.UpdateProductPurchaseStatus(ctx, *item.ProductPurchaseID, models.ProductPurchaseStatusCancelled, nil)
			if err != nil {
				return err
			}
		case models.PurchaseItemTypeProgramOccurrence:
			if item.ParticipationID == nil {
				telemetry.Log().Warn("webhook: program item without participation",
					zap.String("purchase_id", purchase.ID.String()))
				continue
			}
			_, err := o.participations.UpdateParticipationStatus(ctx, *item.ParticipationID, models.ProgramParticipationStatusCancelled)
			if err != nil {
				return err
			}
		}
	}

	o.notify(ctx, purchase, items, models.PurchaseStatusCancelled)
	return nil
}
