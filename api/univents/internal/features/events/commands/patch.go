package commands

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"
	"univents/internal/authz"
	"univents/models"
)

func (c *Commands) Patch(ctx context.Context, payload models.PatchEventInput) (*models.Event, error) {
	ctx, span := telemetry.StartSpan(ctx, "EventService.Patch")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	existing, err := c.events.GetByID(ctx, payload.EventID)
	if err != nil {
		return nil, err
	}

	err = authz.Service.CheckEvent(ctx, ident.Sub.ID, existing.ID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	event := &models.Event{
		FullName:     payload.FullName,
		Acronym:      payload.Acronym,
		Slug:         payload.Slug,
		Description:  payload.Description,
		LogoURL:      payload.LogoURL,
		BannerURL:    payload.BannerURL,
		ContactEmail: payload.ContactEmail,
	}

	return c.events.Patch(ctx, payload.EventID, event)
}
