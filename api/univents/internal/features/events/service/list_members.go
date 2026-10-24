package service

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) ListMembers(ctx context.Context, eventID uuid.UUID) ([]models.EventMember, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListMembers")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	event, err := o.Events.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	err = o.Authz.CheckEvent(ctx, ident.Sub.ID, event.ID, models.EventMemberRoleStaff)
	if err != nil {
		return nil, err
	}

	members, err := o.Events.ListEventMembers(ctx, eventID)
	if err != nil {
		return nil, err
	}

	return members, nil
}
