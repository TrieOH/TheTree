package webhooks

import (
	"context"

	"lib/telemetry"
	"univents/models"

	"go.uber.org/zap"
)

// refund handles a payment.refunded webhook (refund plan slice 3): in a tx,
// guarded approved→refunded (a duplicate refund delivery no-ops), materialized
// rows flipped to cancelled, emitted badges revoked, and stock NOTIFYed (the
// availability ledger only counts pending/approved, so leaving approved frees
// the items automatically). Non-approved purchases are left alone and
// acknowledged.
func (o *Operations) refund(ctx context.Context, purchase *models.Purchase) error {
	if purchase.Status != models.PurchaseStatusApproved {
		telemetry.Log().Info("webhook: refund for non-approved purchase ignored",
			zap.String("purchase_id", purchase.ID.String()),
			zap.String("status", string(purchase.Status)))
		return nil
	}
	return o.tx.WithinTx(ctx, func(ctx context.Context) error {
		updated, err := o.purchases.UpdateStatusIf(ctx, purchase.ID,
			models.PurchaseStatusApproved, models.PurchaseStatusRefunded, nil)
		if err != nil {
			return err
		}
		if updated == nil {
			return nil // guard missed — already flipped by a concurrent delivery
		}
		return o.flipRefunded(ctx, updated)
	})
}

// flipRefunded materializes the refund (slice 3), inside the same tx:
// registrations/product_purchases/participations flipped to cancelled and the
// ticket badges revoked (a refunded buyer must not keep a badge), then NOTIFY
// stock deltas (item ids) so the SSE relay re-reads availability. The D9
// purchase event is deliberately NOT published — the buyer's WS socket closed
// on approval and the hub's framesFor default logs unknown statuses; there is
// no subscriber to a refunded purchase frame.
func (o *Operations) flipRefunded(ctx context.Context, purchase *models.Purchase) error {
	items, err := o.purchases.ListItemsByPurchase(ctx, purchase.ID)
	if err != nil {
		return err
	}

	for _, item := range items {
		switch item.ItemType {
		case models.PurchaseItemTypeTicket:
			if item.RegistrationID == nil {
				telemetry.Log().Warn("webhook: refunded ticket item without registration",
					zap.String("purchase_id", purchase.ID.String()))
				continue
			}
			_, err := o.registrations.UpdateStatus(ctx, *item.RegistrationID, models.RegistrationStatusCancelled, nil)
			if err != nil {
				return err
			}
			err = o.badges.RevokeForRegistration(ctx, *item.RegistrationID, "refunded")
			if err != nil {
				return err
			}
		case models.PurchaseItemTypeProduct:
			if item.ProductPurchaseID == nil {
				telemetry.Log().Warn("webhook: refunded product item without product_purchase",
					zap.String("purchase_id", purchase.ID.String()))
				continue
			}
			_, err := o.productPurchases.UpdateProductPurchaseStatus(ctx, *item.ProductPurchaseID, models.ProductPurchaseStatusCancelled, nil)
			if err != nil {
				return err
			}
		case models.PurchaseItemTypeProgramOccurrence:
			if item.ParticipationID == nil {
				telemetry.Log().Warn("webhook: refunded program item without participation",
					zap.String("purchase_id", purchase.ID.String()))
				continue
			}
			_, err := o.participations.UpdateParticipationStatus(ctx, *item.ParticipationID, models.ProgramParticipationStatusCancelled)
			if err != nil {
				return err
			}
		}
	}

	o.notifyStockOnly(ctx, purchase, items)
	return nil
}

// notifyStockOnly publishes the refunded purchase's stock deltas (D10) —
// item ids only, the SSE relay recomputes from the DB. Unlike approve/cancel
// there is no purchase event: no subscriber exists for a refunded frame.
func (o *Operations) notifyStockOnly(ctx context.Context, purchase *models.Purchase, items []models.PurchaseItem) {
	stock := stockNotification(purchase, items)
	o.publish(ctx, stock)
}
