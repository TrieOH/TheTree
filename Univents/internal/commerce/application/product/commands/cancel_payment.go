package commands

import (
	"context"
	"fmt"
	"log"
	"univents/internal/commerce/domain"
	"univents/internal/shared/sockets"
)

func (uc *CommandService) CancelPayment(ctx context.Context, paymentIntentID string) error {
	purchase, err := uc.purchases.GetByPaymentID(ctx, paymentIntentID)
	if err != nil {
		return fmt.Errorf("failed to fetch purchase for intent %s: %w", paymentIntentID, err)
	}

	switch purchase.Status {
	case domain.PurchaseStatusCancelled:
		log.Printf("[cancel] purchase already cancelled for intent %s", paymentIntentID)
		return nil
	case domain.PurchaseStatusCompleted:
		log.Printf("[cancel] WARNING: received cancel for completed purchase %s intent %s", purchase.ID, paymentIntentID)
		return nil
	}

	if err := uc.purchases.CancelPurchase(ctx, paymentIntentID); err != nil {
		return fmt.Errorf("failed to cancel purchase for intent %s: %w", paymentIntentID, err)
	}

	if purchase.SessionID == nil {
		return nil
	}

	sessionID := *purchase.SessionID
	if err := uc.ws.Notify(sessionID.String(), sockets.WSMessage{
		Type:    "payment_failed",
		Payload: map[string]string{"payment_intent_id": paymentIntentID},
	}); err != nil {
		log.Printf("[cancel] ws already closed for sessions %s: %v", sessionID, err)
	}
	uc.ws.Remove(sessionID.String())

	return nil
}
