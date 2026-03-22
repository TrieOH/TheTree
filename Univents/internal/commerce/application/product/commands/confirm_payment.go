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

func (uc *CommandService) ConfirmPayment(ctx context.Context, payload *paymentsSDK.WebhookPayload) error {
	paymentIntentID := payload.IntentID

	if payload.MercadoPagoData.PaymentMethodID == "pix" && payload.MercadoPagoData.PaymentMethodType == "bank_transfer" {
		var meta struct {
			SessionID uuid.UUID `json:"session_id"`
		}
		if err := json.Unmarshal(payload.Metadata, &meta); err != nil || meta.SessionID == uuid.Nil {
			return fmt.Errorf("missing session_id in metadata for pix intent %s", paymentIntentID)
		}
		if err := uc.ws.Notify(meta.SessionID.String(), sockets.WSMessage{
			Type:    "order_confirmed",
			Payload: map[string]string{"payment_intent_id": paymentIntentID},
		}); err != nil {
			log.Printf("[confirm] ws already closed for pix session %s: %v", meta.SessionID, err)
		}
		return nil
	}

	return uc.finalizeConfirmedPurchase(ctx, paymentIntentID)
}

func (uc *CommandService) finalizeConfirmedPurchase(ctx context.Context, paymentIntentID string) error {
	purchase, err := uc.purchases.GetByPaymentID(ctx, paymentIntentID)
	if err != nil {
		return fmt.Errorf("failed to fetch purchase for intent %s: %w", paymentIntentID, err)
	}

	switch purchase.Status {
	case domain.PurchaseStatusCompleted:
		log.Printf("[confirm] purchase already completed for intent %s", paymentIntentID)
		return nil
	case domain.PurchaseStatusCancelled:
		log.Printf("[confirm] purchase cancelled, ignoring success webhook for %s", paymentIntentID)
		return nil
	}

	// 1. confirm purchase + clean reservation in one TX
	if err := uc.tx.WithinTx(ctx, func(ctx context.Context) error {
		if purchase.SessionID != nil {
			if err := uc.products.DeleteReservation(ctx, *purchase.SessionID); err != nil {
				return err
			}
		}
		if err := uc.purchases.ConfirmPurchase(ctx, paymentIntentID); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to confirm purchase tx for intent %s: %w", paymentIntentID, err)
	}

	// 2. cancel expiry task
	if purchase.SessionID != nil {
		taskID := fmt.Sprintf("%s:%s", *purchase.SessionID, domain.TypeReservationExpired)
		if err := uc.inspector.DeleteTask("default", taskID); err != nil {
			log.Printf("[confirm] could not delete asynq task %s: %v", taskID, err)
		}
	}

	// 3. notify WS if still alive
	if purchase.SessionID != nil {
		sessionID := *purchase.SessionID
		if err := uc.ws.Notify(sessionID.String(), sockets.WSMessage{
			Type:    "order_confirmed",
			Payload: map[string]string{"purchase_id": purchase.ID.String()},
		}); err != nil {
			log.Printf("[confirm] ws already closed for session %s: %v", sessionID, err)
		}
		uc.ws.Remove(sessionID.String())
	}

	// 4. grant ticket permissions
	items, err := uc.purchases.GetTicketIDsByPaymentIntent(ctx, paymentIntentID)
	if err != nil {
		log.Printf("[confirm] failed to fetch ticket ids for %s: %v", paymentIntentID, err)
		return nil
	}

	if len(items) > 0 {
		grants := make([]domain.TicketGrant, 0, len(items))
		for _, item := range items {
			grants = append(grants, domain.TicketGrant{
				TicketID: item.TicketID,
				UserID:   item.UserID,
			})
		}
		task, err := domain.NewGrantTicketPermissionsTask(grants, paymentIntentID)
		if err != nil {
			log.Printf("[confirm] failed to create grant permissions task: %v", err)
			return nil
		}
		if _, err = uc.asynq.EnqueueContext(ctx, task); err != nil {
			log.Printf("[confirm] failed to enqueue grant permissions task: %v", err)
		}
	}

	return nil
}
