package commands

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"
	"univents/models"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (c *Commands) AddMember(ctx context.Context, eventID uuid.UUID, payload models.AddEventMemberRequest) (*models.EventMember, error) {
	ctx, span := telemetry.StartSpan(ctx, "AddMember")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	event, err := c.events.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	if event.OwnerID != ident.Sub.ID {
		member, err := c.events.GetMember(ctx, event.ID, ident.Sub.ID)
		if err != nil {
			return nil, err
		}
		if member.Role != models.EventMemberRoleAdmin {
			return nil, fun.ErrForbidden("insufficient permissions")
		}
	}

	actor, err := c.idx.Actors.GetByEmail(ctx, payload.Email)
	if err != nil {
		return nil, err
	}

	return c.events.AddEventMember(ctx, eventID, actor.ID, payload.Role)
}
