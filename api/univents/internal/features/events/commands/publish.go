package commands

import (
	"context"
	"univents/contracts"
	"univents/internal/shared/errx"

	"github.com/google/uuid"
)

func (uc *Commands) PublishEvent(ctx context.Context, eventID uuid.UUID) error {
	ctx, span := uc.tracer.Start(ctx, "EventService.PublishEvent")
	defer span.End()

	event, err := uc.events.GetByID(ctx, eventID)
	if err != nil {
		return err
	}

	if event.Status != contracts.StatusDraft {
		return errx.Invalid("event").SetMessage("cannot publish non draft event")
	}

	err = uc.events.Publish(ctx, eventID)
	if err != nil {
		return err
	}

	return nil
}
