package commands

import (
	"context"

	"github.com/google/uuid"
)

func (uc *CommandService) Dispatch(ctx context.Context, provider, intentID, event string, eventID uuid.UUID) error {
	return uc.HandleProviderWebhook(ctx, eventID, provider, intentID, event)
}
