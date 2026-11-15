package service

import (
	"context"
	"lib/telemetry"
	"log/slog"
	"univents/models"
	"uuid"
)

func (o *Operations) AddMember(ctx context.Context, eventID uuid.UUID, payload models.AddEventMemberInput) (*models.EventMember, error) {
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
			slog.String("event_id", eventID.String()),
			slog.String("user_id", actor.ID.String()),
			telemetry.Err(err))
	}

	return member, nil
}
