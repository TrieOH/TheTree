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

	ga := uc.gaClient

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return err
	}

	event, err := uc.events.GetEventByID(ctx, eventID)
	if err != nil {
		return err
	}

	allowed, err := ga.Authz.Check().User(sub.ID).
		Object("events").
		Action("publish").
		Scope(event.GoauthScopeID).
		Allowed(ctx)
	if err != nil {
		return err
	}
	if !allowed {
		return errx.Forbidden("event").SetMessage("insufficient permissions")
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
