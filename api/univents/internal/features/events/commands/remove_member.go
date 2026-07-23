package commands

import (
	"context"
	idx "sdk/identityx"
	"univents/models"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (c *Commands) RemoveMember(ctx context.Context, eventID uuid.UUID, payload models.RemoveMemberRequest) error {
	ctx, span := c.tracer.Start(ctx, "RemoveMember")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	event, err := c.events.GetByID(ctx, eventID)
	if err != nil {
		return err
	}

	if event.OwnerID != ident.Sub.ID {
		member, err := c.events.GetMember(ctx, event.ID, ident.Sub.ID)
		if err != nil {
			return err
		}
		if member.Role != models.EventMemberRoleAdmin {
			return fun.ErrForbidden("insufficient permissions")
		}
	}

	actor, err := c.idx.Actors.GetByEmail(ctx, payload.Email)
	if err != nil {
		return err
	}

	return c.events.RemoveEventMember(ctx, eventID, actor.ID)
}
