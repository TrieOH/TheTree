package events

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"
	"univents/models"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (o *Operations) Discontinue(ctx context.Context, eventID uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "EventService.Discontinue")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	event, err := o.events.GetByID(ctx, eventID)
	if err != nil {
		return err
	}

	if event.Status != models.EventStatusActive {
		return fun.ErrBadRequest("cannot discontinue non active event")
	}

	err = o.authz.CheckEvent(ctx, ident.Sub.ID, event.ID, models.EventMemberRoleAdmin)
	if err != nil {
		return err
	}

	return o.events.Discontinue(ctx, eventID)
}
