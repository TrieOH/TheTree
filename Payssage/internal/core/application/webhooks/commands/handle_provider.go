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
		intent, err = uc.intents.Confirm(ctx, id)
	case domain.EventPaymentFailed:
		intent, err = uc.intents.Fail(ctx, id)
	case domain.EventPaymentCancelled:
		intent, err = uc.intents.Cancel(ctx, id)
	default:
		return errx.Invalid("event").SetMessage("unknown event type: " + event)
	}
	if err != nil {
		return err
	}

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
		}
	}

	return nil
}
