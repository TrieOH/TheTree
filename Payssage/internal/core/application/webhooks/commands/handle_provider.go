package commands

import (
	"TriePayments/internal/core/domain"
	"TriePayments/internal/shared/errx"
	"context"
	"encoding/json"
	"log"

	"github.com/google/uuid"
)

func (uc *CommandService) HandleProviderWebhook(ctx context.Context, provider, intentID string, event string) error {
	ctx, span := uc.tracer.Start(ctx, "CommandService.HandleProviderWebhook")
	defer span.End()

	id, err := uuid.Parse(intentID)
	if err != nil {
		return errx.Invalid("intent").SetMessage("invalid intent_id")
	}

	intent, err := uc.intents.GetByID(ctx, id)
	if err != nil {
		return err
	}

	switch event {
	case domain.EventPaymentSucceeded:
		log.Printf("[webhook] confirming intent=%s", id)
		intent, err = uc.intents.Confirm(ctx, id)
	case domain.EventPaymentFailed:
		log.Printf("[webhook] failing intent=%s", id)
		intent, err = uc.intents.Fail(ctx, id)
	case domain.EventPaymentCancelled:
		log.Printf("[webhook] cancelling intent=%s", id)
		intent, err = uc.intents.Cancel(ctx, id)
	default:
		return errx.Invalid("event").SetMessage("unknown event type: " + event)
	}
	if err != nil {
		if errx.IsKind(err, "not_found") {
			// intent already updated by PayIntent, fetch current state and continue
			log.Printf("[webhook] intent=%s already updated, fetching current state", id)
			intent, err = uc.intents.GetByID(ctx, id)
			if err != nil {
				return err
			}
		} else {
			log.Printf("[webhook] failed to update intent=%s event=%s err=%v", id, event, err)
			return err
		}
	}

	log.Printf("[webhook] intent=%s updated to status=%s", intent.ID, intent.Status)

	// build normalized payload
	payloadBytes, err := json.Marshal(domain.WebhookPayload{
		Event:       event,
		IntentID:    intent.ID,
		WorkspaceID: intent.WorkspaceID,
		Amount:      intent.Amount,
		Currency:    intent.Currency,
		Metadata:    intent.Metadata,
	})
	if err != nil {
		return err
	}

	// fetch all registered endpoints for this workspace
	endpoints, err := uc.endpoints.ListByWorkspace(ctx, intent.WorkspaceID)
	if err != nil {
		return err
	}

	log.Printf("[webhook] found %d endpoints for workspace=%s", len(endpoints), intent.WorkspaceID)

	// enqueue delivery task per endpoint
	for _, endpoint := range endpoints {
		delivery, err := domain.NewWebhookDelivery(endpoint.ID, intent.ID, event, payloadBytes)
		if err != nil {
			log.Printf("[webhook] failed to create delivery object for endpoint %s: %v", endpoint.ID, err)
			continue
		}
		created, err := uc.deliveries.Create(ctx, *delivery)
		if err != nil {
			log.Printf("[webhook] failed to create delivery record for endpoint %s: %v", endpoint.ID, err)
			continue
		}

		task, err := domain.NewDeliverWebhookTask(created.ID, endpoint.ID, endpoint.URL, endpoint.Secret, payloadBytes)
		if err != nil {
			log.Printf("[webhook] failed to create delivery task for endpoint %s: %v", endpoint.ID, err)
			continue
		}

		if _, err = uc.asynq.EnqueueContext(context.Background(), task); err != nil {
			log.Printf("[webhook] failed to enqueue delivery task for endpoint %s: %v", endpoint.ID, err)
		} else {
			log.Printf("[webhook] enqueued delivery for endpoint=%s url=%s", endpoint.ID, endpoint.URL)
		}
	}

	return nil
}
