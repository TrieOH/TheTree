package service

import (
	"context"
	"lib/telemetry"
	"log/slog"
	"univents/models"
	"uuid"
)

func (o *Operations) RemoveMember(ctx context.Context, eventID uuid.UUID, payload models.RemoveMemberInput) error {
	event, err := o.Events.GetByID(ctx, eventID)
	if err != nil {
		return err
	}

	_, err = o.Authz.Event(ctx, event.ID, models.EventMemberRoleAdmin)
	if err != nil {
		return err
	}

	actor, err := o.IDX.Actors.GetByEmail(ctx, payload.Email)
	if err != nil {
		return err
	}

	err = o.Events.RemoveEventMember(ctx, eventID, actor.ID)
	if err != nil {
		return err
	}

	err = o.Badges.RevokeStaffBadgesForUser(ctx, eventID, actor.ID)
	if err != nil {
		telemetry.Log().Error("failed to revoke staff badges",
			slog.String("event_id", eventID.String()),
			slog.String("user_id", actor.ID.String()),
			telemetry.Err(err))
	}

	return nil
}
