package commands

import (
	"TriePayments/internal/core/domain"
	"context"
	"encoding/json"

	"github.com/google/uuid"
)

func (uc *CommandService) CreateWebhookEvent(ctx context.Context, provider, eventType string, payload json.RawMessage) (uuid.UUID, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, err
	}

	event := domain.WebhookEvent{
		ID:        id,
		Provider:  provider,
		EventType: eventType,
		Payload:   payload,
	}
	created, err := uc.events.Create(ctx, event)
	if err != nil {
		return uuid.Nil, err
	}
	return created.ID, nil
}
