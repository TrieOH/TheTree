package service

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) ListDraft(ctx context.Context, eventID uuid.UUID) ([]models.Edition, error) {
	ctx, span := telemetry.StartSpan(ctx, "EditionService.ListDraft")
	defer span.End()

	event, err := o.Events.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Event(ctx, event.ID, models.EventMemberRoleStaff)
	if err != nil {
		return nil, err
	}

	return o.Editions.ListDraft(ctx, event.ID)
}
