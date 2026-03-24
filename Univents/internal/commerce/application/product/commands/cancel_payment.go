package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"univents/internal/shared/sockets"

	paymentsSDK "github.com/TrieOH/TriePaymentsSDK"
	"github.com/google/uuid"
)

func (uc *CommandService) CancelPayment(ctx context.Context, payload *paymentsSDK.WebhookPayload) error {
	paymentIntentID := payload.IntentID

	var meta struct {
		SessionID uuid.UUID `json:"session_id"`
		UserID    uuid.UUID `json:"user_id"`
	}
	if err := json.Unmarshal(payload.Metadata, &meta); err != nil || meta.SessionID == uuid.Nil || meta.UserID == uuid.Nil {
		return fmt.Errorf("missing session_id or user_id in metadata for intent %s", paymentIntentID)
	}

	if err := uc.sessions.Delete(ctx, meta.UserID, meta.SessionID); err != nil {
		log.Printf("[cancel] failed to delete session %s: %v", meta.SessionID, err)
	}

	if err := uc.ws.Notify(meta.SessionID.String(), sockets.WSMessage{
		Type:    "payment_failed",
		Payload: map[string]string{"payment_intent_id": paymentIntentID},
	}); err != nil {
		log.Printf("[cancel] ws already closed for session %s: %v", meta.SessionID, err)
	}
	uc.ws.Remove(meta.SessionID.String())

	return nil
}
