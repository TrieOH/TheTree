package events

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"
	"univents/internal/authz"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) RemoveMember(ctx context.Context, eventID uuid.UUID, payload models.RemoveMemberRequest) error {
	ctx, span := telemetry.StartSpan(ctx, "RemoveMember")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	event, err := o.events.GetByID(ctx, eventID)
	if err != nil {
		return err
	}

	err = authz.Service.CheckEvent(ctx, ident.Sub.ID, event.ID, models.EventMemberRoleAdmin)
	if err != nil {
		return err
	}

	actor, err := o.idx.Actors.GetByEmail(ctx, payload.Email)
	if err != nil {
		return err
	}

	return o.events.RemoveEventMember(ctx, eventID, actor.ID)
}
