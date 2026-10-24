package service

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) ListMembers(ctx context.Context, eventID uuid.UUID) ([]models.EventMember, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListMembers")
	defer span.End()

	event, err := o.Events.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Event(ctx, event.ID, models.EventMemberRoleStaff)
	if err != nil {
		return nil, err
	}

	members, err := o.Events.ListEventMembers(ctx, eventID)
	if err != nil {
		return nil, err
	}

	return members, nil
}
