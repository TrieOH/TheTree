package commands

import (
	"TriePayments/internal/core/domain"
	"TriePayments/internal/shared/errx"
	"context"
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

func (uc *CommandService) HandleProviderWebhook(ctx context.Context, eventID uuid.UUID, provider, intentID string, event string) (err error) {
	ctx, span := uc.tracer.Start(ctx, "CommandService.HandleProviderWebhook")
	defer span.End()

	var id uuid.UUID
	id, err = uuid.Parse(intentID)
	if err != nil {
		return errx.Invalid("intent").SetMessage("invalid intent_id")
	}

	var intent *domain.Intent
	intent, err = uc.intents.GetByID(ctx, id)
	if err != nil {
		return err
	}

	switch event {
	case domain.EventPaymentSucceeded:
		if !alreadyInTargetState(event, intent.Status) {
			log.Printf("[webhook] confirming intent=%s", id)
			intent, err = uc.intents.Confirm(ctx, id)
		} else {
			log.Printf("[webhook] intent=%s already in target state, skipping mutation", id)
		}
	case domain.EventPaymentFailed:
		if !alreadyInTargetState(event, intent.Status) {
			log.Printf("[webhook] failing intent=%s", id)
			intent, err = uc.intents.Fail(ctx, id)
		} else {
			log.Printf("[webhook] intent=%s already in target state, skipping mutation", id)
		}
	case domain.EventPaymentCancelled:
		if !alreadyInTargetState(event, intent.Status) {
			log.Printf("[webhook] cancelling intent=%s", id)
			intent, err = uc.intents.Cancel(ctx, id)
		} else {
			log.Printf("[webhook] intent=%s already in target state, skipping mutation", id)
		}
	default:
		return errx.Invalid("event").SetMessage("unknown event type: " + event)
	}
	if err != nil {
		log.Printf("[webhook] failed to update intent=%s event=%s err=%v", id, event, err)
		return err
	}

	log.Printf("[webhook] intent=%s status=%s", intent.ID, intent.Status)

	// build normalized payload
	var payloadBytes []byte
	payloadBytes, err = json.Marshal(domain.WebhookPayload{
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
	var endpoints []domain.WebhookEndpoint
	endpoints, err = uc.endpoints.ListByWorkspace(ctx, intent.WorkspaceID)
	if err != nil {
		return err
	}

	log.Printf("[webhook] found %d endpoints for workspace=%s", len(endpoints), intent.WorkspaceID)

	// enqueue delivery task per endpoint
	for _, endpoint := range endpoints {
		var delivery *domain.WebhookDelivery
		delivery, err = domain.NewWebhookDelivery(endpoint.ID, intent.ID, event, payloadBytes)
		if err != nil {
			log.Printf("[webhook] failed to create delivery object for endpoint %s: %v", endpoint.ID, err)
			continue
		}
		var created *domain.WebhookDelivery
		created, err = uc.deliveries.Create(ctx, *delivery)
		if err != nil {
			log.Printf("[webhook] failed to create delivery record for endpoint %s: %v", endpoint.ID, err)
			continue
		}

		var task *asynq.Task
		task, err = domain.NewDeliverWebhookTask(created.ID, endpoint.ID, endpoint.URL, endpoint.Secret, payloadBytes)
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

	if eventID != uuid.Nil {
		uc.EnrichWebhookEvent(ctx, eventID, intent.WorkspaceID, intent.ID, intentID)
	}

	return nil
}

func alreadyInTargetState(event string, status domain.IntentStatus) bool {
	switch event {
	case domain.EventPaymentSucceeded:
		return status == domain.IntentStatusSucceeded
	case domain.EventPaymentFailed:
		return status == domain.IntentStatusFailed
	case domain.EventPaymentCancelled:
		return status == domain.IntentStatusCancelled
	}
	return false
}
