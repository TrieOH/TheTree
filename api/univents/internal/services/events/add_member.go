package events

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"
	"univents/models"

	"github.com/google/uuid"

	"go.uber.org/zap"
)

func (o *Operations) AddMember(ctx context.Context, eventID uuid.UUID, payload models.AddEventMemberInput) (*models.EventMember, error) {
	ctx, span := telemetry.StartSpan(ctx, "AddMember")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	event, err := o.events.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	err = o.authz.CheckEvent(ctx, ident.Sub.ID, event.ID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	actor, err := o.idx.Actors.GetByEmail(ctx, payload.Email)
	if err != nil {
		return nil, err
	}

	member, err := o.events.AddEventMember(ctx, eventID, actor.ID, payload.Role)
	if err != nil {
		return nil, err
	}

	err = o.badges.AwardStaffBadgesForUser(ctx, eventID, actor.ID)
	if err != nil {
		telemetry.Log().Error("failed to award staff badges",
			zap.String("event_id", eventID.String()),
			zap.String("user_id", actor.ID.String()),
			zap.Error(err))
	}

	return member, nil
}
