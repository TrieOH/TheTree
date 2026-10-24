package service

import (
	"context"
	idx "sdk/identityx"
	"univents/models"
)

func (o *Operations) Create(ctx context.Context, payload models.CreateEventInput) (*models.Event, error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	event := &models.Event{
		OwnerID:      ident.Sub.ID,
		FullName:     payload.FullName,
		Acronym:      payload.Acronym,
		Slug:         payload.Slug,
		Description:  payload.Description,
		Status:       models.EventStatusDraft,
		ContactEmail: payload.ContactEmail,
	}

	var created *models.Event
	err = o.Tx.WithinTx(ctx, func(ctx context.Context) error {
		created, err = o.Events.Create(ctx, event)
		if err != nil {
			return err
		}

		_, err = o.Events.AddEventMember(ctx, created.ID, ident.Sub.ID, models.EventMemberRoleOwner)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return created, nil
}
