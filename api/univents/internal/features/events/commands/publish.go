package commands

import (
	"context"
	idx "sdk/identityx"
	"univents/models"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (c *Commands) Publish(ctx context.Context, eventID uuid.UUID) error {
	ctx, span := c.tracer.Start(ctx, "EventService.Publish")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	event, err := c.events.GetByID(ctx, eventID)
	if err != nil {
		return err
	}

	if event.Status != models.EventStatusDraft {
		return fun.ErrBadRequest("cannot publish non draft event")
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

	return c.events.Publish(ctx, eventID)
}
