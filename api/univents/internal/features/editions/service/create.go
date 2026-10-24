package service

import (
	"context"
	"lib/telemetry"
	"univents/models"
)

func (o *Operations) Create(ctx context.Context, payload models.CreateEditionInput) (*models.Edition, error) {
	ctx, span := telemetry.StartSpan(ctx, "EditionService.Create")
	defer span.End()

	event, err := o.Events.GetByID(ctx, payload.EventID)
	if err != nil {
		return nil, err
	}

	ident, err := o.Authz.Event(ctx, event.ID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	edition := &models.Edition{
		EventID:   payload.EventID,
		Name:      payload.Name,
		Slug:      payload.Slug,
		IsDraft:   true,
		StartsAt:  payload.StartsAt,
		EndsAt:    payload.EndsAt,
		CreatedBy: ident.Sub.ID,
	}

	return o.Editions.Create(ctx, edition)
}
