package commands

import (
	"context"
	"errors"
	"fmt"
	"log"
	"univents/internal/commerce/domain"
	"univents/internal/plataform/telemetry"
	"univents/internal/shared/sockets"

	"github.com/hibiken/asynq"
)

func (uc *CommandService) ConfirmPayment(ctx context.Context, paymentIntentID string) error {
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

	msg := fmt.Sprintf("malformed sessionID for purchase %s", paymentIntentID)
	if purchase.SessionID == nil {
		telemetry.Log().Error(msg)
		return errors.New(msg)
	}
	sessionID := *purchase.SessionID

	// 1. mark items as sold in db
	if err := uc.tx.WithinTx(ctx, func(ctx context.Context) error {
		if err := uc.products.DeleteReservation(ctx, sessionID); err != nil {
			return err
		}
		if err := uc.purchases.ConfirmPurchase(ctx, paymentIntentID); err != nil {
			return err
		}
		return nil
	}); err != nil {
		if err := uc.ws.Notify(sessionID.String(), sockets.WSMessage{
			Type:    "order_failed",
			Payload: map[string]string{"payment_intent_id": paymentIntentID},
		}); err != nil {
			log.Printf("[confirm] ws already closed for session %s: %v", sessionID, err)
		}
		uc.ws.Remove(sessionID.String())
		return nil
	}

	// 2. cancel the asynq expiry task so it doesn't fire after successful payment
	taskID := fmt.Sprintf("%s:%s:%s", sessionID, paymentIntentID, domain.TypeReservationExpired)
	if err := uc.inspector.DeleteTask("default", taskID); err != nil {
		// task may have already fired or doesn't exist — log but don't fail
		log.Printf("[confirm] could not delete asynq task %s: %v", taskID, err)
	}

	// 3. notify the open purchase socket
	if err := uc.ws.Notify(sessionID.String(), sockets.WSMessage{
		Type:    "order_confirmed",
		Payload: map[string]string{"purchase_id": purchase.ID.String()},
	}); err != nil {
		log.Printf("[confirm] ws already closed for session %s: %v", sessionID, err)
	}

	uc.ws.Remove(sessionID.String())

	items, err := uc.purchases.GetTicketIDsByPaymentIntent(ctx, paymentIntentID)
	if err != nil {
		log.Printf("[confirm] failed to fetch ticket ids for %s: %v", paymentIntentID, err)
	} else if len(items) > 0 {
		grants := make([]domain.TicketGrant, 0, len(items))
		for _, item := range items {
			grants = append(grants, domain.TicketGrant{
				TicketID: item.TicketID,
				UserID:   item.UserID,
			})
		}

		var task *asynq.Task
		task, err = domain.NewGrantTicketPermissionsTask(grants, paymentIntentID)
		if err != nil {
			log.Printf("[confirm] failed to create grant permissions task: %v", err)
		} else if _, err = uc.asynq.EnqueueContext(ctx, task); err != nil {
			log.Printf("[confirm] failed to enqueue grant permissions task: %v", err)
		}
	}

	return nil
}
