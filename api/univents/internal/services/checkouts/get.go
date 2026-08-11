package checkouts

import (
	"context"
	"time"

	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Get is the resume read: the purchase owner-scoped (the caller must be the
// purchase's purchaser_id — anything else is NOT_FOUND via the repo, no
// existence leak), its items from the availability ledger, and — only when
// the purchase is still pending and carries an intent — the intent's
// current status. Never transitions state.
func (o *Operations) Get(ctx context.Context, purchaseID, purchaserID uuid.UUID) (*Resume, error) {
	ctx, span := telemetry.StartSpan(ctx, "CheckoutsService.Get")
	defer span.End()

	purchase, err := o.purchases.GetByIDForOwner(ctx, purchaseID, purchaserID)
	if err != nil {
		return nil, err
	}
	items, err := o.purchases.ListItemsByPurchase(ctx, purchaseID)
	if err != nil {
		return nil, err
	}

	resume := &Resume{Purchase: *purchase, Items: items}
	if purchase.Status == models.PurchaseStatusPending && purchase.PayssageIntentID != nil {
		resume.IntentStatus = o.intentStatus(ctx, *purchase.PayssageIntentID)
	}
	return resume, nil
}

// intentStatus fetches the intent's status with the retry budget in
// Operations (up to 3 tries by default), then degrades gracefully: on
// persistent failure it logs loudly — a resume that cannot reach Payssage
// still returns the purchase's own state, because the purchases table is
// the source of truth and `intent_status` is only a supplement.
func (o *Operations) intentStatus(ctx context.Context, intentID uuid.UUID) *string {
	var lastErr error
	for attempt := 1; attempt <= o.intentAttempts; attempt++ {
		intent, err := o.intents.GetIntent(ctx, intentID)
		if err == nil {
			return new(string(intent.Status))
		}
		lastErr = err
		if attempt < o.intentAttempts {
			select {
			case <-ctx.Done():
				telemetry.Log().Warn("checkouts: intent status fetch cancelled",
					zap.String("intent_id", intentID.String()),
					zap.Error(ctx.Err()))
				return nil
			case <-time.After(o.intentRetryDelay):
			}
		}
	}
	telemetry.Log().Error("checkouts: intent status unavailable after retries — degrading (purchases table remains the source of truth)",
		zap.String("intent_id", intentID.String()),
		zap.Int("attempts", o.intentAttempts),
		zap.Error(lastErr))
	return nil
}
