package events

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"
	"univents/models"

	"github.com/google/uuid"

	"go.uber.org/zap"
)

func (o *Operations) RemoveMember(ctx context.Context, eventID uuid.UUID, payload models.RemoveMemberInput) error {
	ctx, span := telemetry.StartSpan(ctx, "RemoveMember")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	event, err := o.events.GetByID(ctx, eventID)
	if err != nil {
		return err
	}

	err = o.authz.CheckEvent(ctx, ident.Sub.ID, event.ID, models.EventMemberRoleAdmin)
	if err != nil {
		return err
	}

	actor, err := o.idx.Actors.GetByEmail(ctx, payload.Email)
	if err != nil {
		return err
	}

	err = o.events.RemoveEventMember(ctx, eventID, actor.ID)
	if err != nil {
		return err
	}

	err = o.badges.RevokeStaffBadgesForUser(ctx, eventID, actor.ID)
	if err != nil {
		telemetry.Log().Error("failed to revoke staff badges",
			zap.String("event_id", eventID.String()),
			zap.String("user_id", actor.ID.String()),
			zap.Error(err))
	}

	return nil
}
