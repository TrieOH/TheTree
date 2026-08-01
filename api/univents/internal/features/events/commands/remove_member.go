package commands

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"
	"univents/internal/authz"
	"univents/models"

	"github.com/google/uuid"
)

func (c *Commands) RemoveMember(ctx context.Context, eventID uuid.UUID, payload models.RemoveMemberRequest) error {
	ctx, span := telemetry.StartSpan(ctx, "RemoveMember")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	event, err := c.events.GetByID(ctx, eventID)
	if err != nil {
		return err
	}

	err = authz.Service.CheckEvent(ctx, ident.Sub.ID, event.ID, models.EventMemberRoleAdmin)
	if err != nil {
		return err
	}

	actor, err := c.idx.Actors.GetByEmail(ctx, payload.Email)
	if err != nil {
		return err
	}

	return c.events.RemoveEventMember(ctx, eventID, actor.ID)
}
