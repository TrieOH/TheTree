package editions

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"
	"univents/internal/authz"
	"univents/models"
)

func (o *Operations) Create(ctx context.Context, payload models.CreateEditionInput) (*models.Edition, error) {
	ctx, span := telemetry.StartSpan(ctx, "EditionService.Create")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	event, err := o.events.GetByID(ctx, payload.EventID)
	if err != nil {
		return nil, err
	}

	err = authz.Service.CheckEvent(ctx, ident.Sub.ID, event.ID, models.EventMemberRoleAdmin)
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

	return o.editions.Create(ctx, edition)
}
