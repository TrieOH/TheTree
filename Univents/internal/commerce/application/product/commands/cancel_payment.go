package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"univents/internal/commerce/domain"
	"univents/internal/shared/sockets"

	paymentsSDK "github.com/TrieOH/TriePaymentsSDK"
	"github.com/google/uuid"
)

func (uc *CommandService) CancelPayment(ctx context.Context, payload *paymentsSDK.WebhookPayload) error {
	paymentIntentID := payload.IntentID

	if payload.MercadoPagoData.PaymentMethodID == "pix" && payload.MercadoPagoData.PaymentMethodType == "bank_transfer" {
		var meta struct {
			SessionID uuid.UUID `json:"session_id"`
		}
		if err := json.Unmarshal(payload.Metadata, &meta); err != nil || meta.SessionID == uuid.Nil {
			return fmt.Errorf("missing session_id in metadata for pix intent %s", paymentIntentID)
		}
		if err := uc.ws.Notify(meta.SessionID.String(), sockets.WSMessage{
			Type:    "payment_failed",
			Payload: map[string]string{"payment_intent_id": paymentIntentID},
		}); err != nil {
			log.Printf("[cancel] ws already closed for pix session %s: %v", meta.SessionID, err)
		}
		uc.ws.Remove(meta.SessionID.String())
		return nil
	}

	return uc.finalizeFailedPurchase(ctx, paymentIntentID)
}

func (uc *CommandService) finalizeFailedPurchase(ctx context.Context, paymentIntentID string) error {
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
		log.Printf("[cancel] ws already closed for session %s: %v", sessionID, err)
	}
	uc.ws.Remove(sessionID.String())

	return nil
}
