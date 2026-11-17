package service

import (
	"context"
	"encoding/json"
	"lib/telemetry"
	"log/slog"
	"univents/internal/notify"
	"univents/models"
)

// notifyStock publishes the reservation's stock deltas (D10): item_ids only
// — the SSE relay re-queries availability from the DB. Fire-and-forget: a
// missed notification is a stale snapshot, never data loss.
func (o *Operations) notifyStock(ctx context.Context, purchase *models.Purchase, items []models.PurchaseItem) {
	stock := notify.StockNotification{Kind: notify.KindStock, EditionID: purchase.EditionID}
	for _, item := range items {
		stock.ItemIDs = append(stock.ItemIDs, item.ItemID)
	}
	o.publish(ctx, stock)
}

// notifyPurchase publishes a purchase event (D9) — the free-order confirm
// and the revert cancel. The WS hub maps the status to the frame.
func (o *Operations) notifyPurchase(ctx context.Context, purchase *models.Purchase, status models.PurchaseStatus) {
	o.publish(ctx, notify.PurchaseNotification{
		Kind:       notify.KindPurchase,
		EditionID:  purchase.EditionID,
		PurchaseID: purchase.ID,
		Status:     string(status),
	})
}

func (o *Operations) publish(ctx context.Context, payload any) {
	raw, err := json.Marshal(payload)
	if err != nil {
		telemetry.Log().Error("checkout: marshal notifier payload",
			slog.String("channel", channelUniventsChanges),
			telemetry.Err(err))
		return
	}
	err = o.Notifier.Notify(ctx, channelUniventsChanges, string(raw))
	if err != nil {
		telemetry.Log().Error("checkout: publish to notifier",
			slog.String("channel", channelUniventsChanges),
			slog.String("payload", string(raw)),
			telemetry.Err(err))
	}
}
