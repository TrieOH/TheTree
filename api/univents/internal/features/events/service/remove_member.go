package service

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"

	"go.uber.org/zap"
)

func (o *Operations) RemoveMember(ctx context.Context, eventID uuid.UUID, payload models.RemoveMemberInput) error {
	ctx, span := telemetry.StartSpan(ctx, "RemoveMember")
	defer span.End()

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
			zap.String("event_id", eventID.String()),
			zap.String("user_id", actor.ID.String()),
			zap.Error(err))
	}

	return nil
}
