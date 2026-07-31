package commands

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"
	"univents/internal/authz"
	"univents/models"

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

	err = authz.Service.CheckEvent(ctx, ident.Sub.ID, event.ID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	actor, err := c.idx.Actors.GetByEmail(ctx, payload.Email)
	if err != nil {
		return nil, err
	}

	return c.events.AddEventMember(ctx, eventID, actor.ID, payload.Role)
}
