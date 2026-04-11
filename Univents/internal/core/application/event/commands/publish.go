package commands

import (
	"context"
	"univents/internal/core/domain"
	"univents/internal/shared/authz"
	"univents/internal/shared/errx"

	"github.com/google/uuid"
)

func (uc *CommandService) PublishEvent(ctx context.Context, eventID uuid.UUID) error {
	ctx, span := uc.tracer.Start(ctx, "EventService.PublishEvent")
	defer span.End()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return err
	}

	event, err := uc.events.GetByID(ctx, eventID)
	if err != nil {
		return err
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("publish"),
		authz.Resource("event", event.ID.String()),
	); err != nil {
		return err
	}

	if event.Status != domain.StatusDraft {
		return errx.Invalid("event").SetMessage("cannot publish non draft event")
	}

	err = uc.events.PublishEvent(ctx, eventID)
	if err != nil {
		return err
	}

	return nil
}
