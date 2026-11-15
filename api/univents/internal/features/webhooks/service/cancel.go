package service

import (
	"context"
	"log/slog"

	"lib/telemetry"
	"univents/models"
)

// cancel handles a terminal failure webhook (failed/rejected/cancelled): in
// a tx, guarded pending→target (a stale failure delivery can't cancel an
// approved purchase), materialized rows flipped to cancelled, NOTIFY (stock
// freed + purchase event). Non-pending purchases are left alone and
// acknowledged.
//
// The target status distinguishes why the purchase died (D4): the customer
// cancelling is "cancelled", a declined payment is "failed" (payment.failed)
// or "rejected" (payment.rejected — e.g. MP's risk engine). The normalized
// reason from the delivery envelope lands in status_reason.
func (o *Operations) cancel(ctx context.Context, purchase *models.Purchase, to models.PurchaseStatus, reason *string) error {
	if purchase.Status != models.PurchaseStatusPending {
		telemetry.Log().Info("webhook: failure for non-pending purchase ignored",
			slog.String("purchase_id", purchase.ID.String()),
			slog.String("status", string(purchase.Status)))
		return nil
	}
	return o.Tx.WithinTx(ctx, func(ctx context.Context) error {
		updated, err := o.Purchases.UpdateStatusIf(ctx, purchase.ID,
			models.PurchaseStatusPending, to, reason)
		if err != nil {
			return err
		}
		if updated == nil {
			return nil // guard missed — already flipped by a concurrent delivery
		}
		return o.flipCancelled(ctx, updated, to)
	})
}

// flipCancelled materializes the cancellation (D4), inside the same tx:
// registrations/product_purchases/participations flipped to cancelled, then
// NOTIFY with the purchase's target status (stock freed + purchase event).
func (o *Operations) flipCancelled(ctx context.Context, purchase *models.Purchase, to models.PurchaseStatus) error {
	items, err := o.Purchases.ListItemsByPurchase(ctx, purchase.ID)
	if err != nil {
		return err
	}

	for _, item := range items {
		switch item.ItemType {
		case models.PurchaseItemTypeTicket:
			if item.RegistrationID == nil {
				telemetry.Log().Warn("webhook: ticket item without registration",
					slog.String("purchase_id", purchase.ID.String()))
				continue
			}
			_, err := o.Registrations.UpdateStatus(ctx, *item.RegistrationID, models.RegistrationStatusCancelled, nil)
			if err != nil {
				return err
			}
		case models.PurchaseItemTypeProduct:
			if item.ProductPurchaseID == nil {
				telemetry.Log().Warn("webhook: product item without product_purchase",
					slog.String("purchase_id", purchase.ID.String()))
				continue
			}
			_, err := o.ProductPurchases.UpdateProductPurchaseStatus(ctx, *item.ProductPurchaseID, models.ProductPurchaseStatusCancelled, nil)
			if err != nil {
				return err
			}
		case models.PurchaseItemTypeProgramOccurrence:
			if item.ParticipationID == nil {
				telemetry.Log().Warn("webhook: program item without participation",
					slog.String("purchase_id", purchase.ID.String()))
				continue
			}
			_, err := o.Participations.UpdateParticipationStatus(ctx, *item.ParticipationID, models.ProgramParticipationStatusCancelled)
			if err != nil {
				return err
			}
		}
	}

	o.notify(ctx, purchase, items, to)
	return nil
}
