package events

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"
	"univents/models"
)

func (o *Operations) Patch(ctx context.Context, payload models.PatchEventInput) (*models.Event, error) {
	ctx, span := telemetry.StartSpan(ctx, "EventService.Patch")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	existing, err := o.events.GetByID(ctx, payload.EventID)
	if err != nil {
		return nil, err
	}

	err = o.authz.CheckEvent(ctx, ident.Sub.ID, existing.ID, models.EventMemberRoleAdmin)
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

	return o.events.Patch(ctx, payload.EventID, event)
}
