package commands

import (
	"context"
	"log"

	"github.com/google/uuid"
)

func (uc *CommandService) EnrichWebhookEvent(ctx context.Context, eventID, workspaceID, intentID uuid.UUID, externalID string) {
	if _, err := uc.events.Enrich(ctx, eventID, workspaceID, intentID, externalID); err != nil {
		log.Printf("[webhook_event] failed to enrich event=%s err=%v", eventID, err)
	}
}
