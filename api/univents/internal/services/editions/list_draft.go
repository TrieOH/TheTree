package editions

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) ListDraft(ctx context.Context, eventID uuid.UUID) ([]models.Edition, error) {
	ctx, span := telemetry.StartSpan(ctx, "EditionService.ListDraft")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	event, err := o.events.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	err = o.authz.CheckEvent(ctx, ident.Sub.ID, event.ID, models.EventMemberRoleStaff)
	if err != nil {
		return nil, err
	}

	return o.editions.ListDraft(ctx, event.ID)
}
