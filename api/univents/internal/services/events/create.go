package events

import (
	"context"
	"lib/database"
	"lib/telemetry"
	idx "sdk/identityx"
	"univents/models"
)

func (o *Operations) Create(ctx context.Context, payload models.CreateEventInput) (*models.Event, error) {
	ctx, span := telemetry.StartSpan(ctx, "Create")
	defer span.End()

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
	err = database.RunTx(ctx, func(ctx context.Context) error {
		created, err = o.events.Create(ctx, event)
		if err != nil {
			return err
		}

		_, err = o.events.AddEventMember(ctx, created.ID, ident.Sub.ID, models.EventMemberRoleOwner)
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
