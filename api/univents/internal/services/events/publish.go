package events

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"
	"univents/models"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (o *Operations) Publish(ctx context.Context, eventID uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "EventService.Publish")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	event, err := o.events.GetByID(ctx, eventID)
	if err != nil {
		return err
	}

	if event.Status != models.EventStatusDraft {
		return fun.ErrBadRequest("cannot publish non draft event")
	}

	err = o.authz.CheckEvent(ctx, ident.Sub.ID, event.ID, models.EventMemberRoleAdmin)
	if err != nil {
		return err
	}

	return o.events.Publish(ctx, eventID)
}
