package commands

import (
	"context"
	idx "sdk/identityx"
	"univents/models"
)

func (c *Commands) Create(ctx context.Context, payload models.CreateEventRequest) (*models.Event, error) {
	ctx, span := c.tracer.Start(ctx, "Create")
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

	return c.events.Create(ctx, event)
}
