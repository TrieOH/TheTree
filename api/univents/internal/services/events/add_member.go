package events

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"
	"univents/internal/authz"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) AddMember(ctx context.Context, eventID uuid.UUID, payload models.AddEventMemberInput) (*models.EventMember, error) {
	ctx, span := telemetry.StartSpan(ctx, "AddMember")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	event, err := o.events.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	err = authz.Service.CheckEvent(ctx, ident.Sub.ID, event.ID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	actor, err := o.idx.Actors.GetByEmail(ctx, payload.Email)
	if err != nil {
		return nil, err
	}

	return o.events.AddEventMember(ctx, eventID, actor.ID, payload.Role)
}
