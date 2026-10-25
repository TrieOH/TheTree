package service

import (
	"context"
	"univents/models"
)

func (o *Operations) Create(ctx context.Context, payload models.CreateEditionInput) (*models.Edition, error) {
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
