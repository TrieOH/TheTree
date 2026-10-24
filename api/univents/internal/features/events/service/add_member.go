package service

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"

	"go.uber.org/zap"
)

func (o *Operations) AddMember(ctx context.Context, eventID uuid.UUID, payload models.AddEventMemberInput) (*models.EventMember, error) {
	ctx, span := telemetry.StartSpan(ctx, "AddMember")
	defer span.End()

	event, err := o.Events.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Event(ctx, event.ID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	actor, err := o.IDX.Actors.GetByEmail(ctx, payload.Email)
	if err != nil {
		return nil, err
	}

	member, err := o.Events.AddEventMember(ctx, eventID, actor.ID, payload.Role)
	if err != nil {
		return nil, err
	}

	err = o.Badges.AwardStaffBadgesForUser(ctx, eventID, actor.ID)
	if err != nil {
		telemetry.Log().Error("failed to award staff badges",
			zap.String("event_id", eventID.String()),
			zap.String("user_id", actor.ID.String()),
			zap.Error(err))
	}

	return member, nil
}
