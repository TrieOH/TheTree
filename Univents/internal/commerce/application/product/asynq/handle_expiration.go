package async

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"univents/internal/commerce/domain"
	"univents/internal/shared/sockets"

	paymentsSDK "github.com/TrieOH/TriePaymentsSDK"
	"github.com/hibiken/asynq"
)

func (uc *AsynqHandlers) HandleProductReservationExpiration(ctx context.Context, t *asynq.Task) error {
	var p domain.ReservationExpiredPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}

	log.Printf("[task] reservation expired for session %s", p.SessionID)

	if err := uc.tx.WithinTx(ctx, func(ctx context.Context) error {
		if err := uc.products.UnreserveItems(ctx, p.SessionID); err != nil {
			return fmt.Errorf("failed to unreserve items: %w", err)
		}
		if err := uc.purchases.CancelPurchase(ctx, p.PaymentIntentID); err != nil {
			return fmt.Errorf("failed to cancel purchase: %w", err)
		}
		return nil
	}); err != nil {
		return err
	}

	// 2. Cancel payment intent
	if _, err := uc.payments.CancelIntent(ctx, p.PaymentIntentID); err != nil {
		if paymentsSDK.IsNotFound(err) {
			// intent already succeeded or was canceled — skip
			log.Printf("[task] intent %s already gone, skipping cancel", p.PaymentIntentID)
		} else {
			return fmt.Errorf("failed to cancel payment intent: %w", err)
		}
	}

	// 3. Notify WS connection if still alive
	if err := uc.ws.Notify(p.SessionID.String(), sockets.WSMessage{
		Type:    "reservation_expired",
		Payload: "your reservation timed out",
	}); err != nil {
		// connection already gone — not an error worth retrying the task for
		log.Printf("[task] ws already closed for session %s: %v", p.SessionID, err)
	}

	uc.ws.Remove(p.SessionID.String())

	return nil
}
