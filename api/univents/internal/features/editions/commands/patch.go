package commands

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"
	"univents/internal/authz"
	"univents/models"
)

func (c *Commands) Patch(ctx context.Context, payload models.PatchEditionInput) (*models.Edition, error) {
	ctx, span := telemetry.StartSpan(ctx, "EditionService.Patch")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	existing, err := c.editions.GetByID(ctx, payload.EditionID)
	if err != nil {
		return nil, err
	}

	event, err := c.events.GetByID(ctx, existing.EventID)
	if err != nil {
		return nil, err
	}

	err = authz.Service.CheckEvent(ctx, ident.Sub.ID, event.ID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	edition := &models.Edition{
		Name:                payload.Name,
		Slug:                payload.Slug,
		Tagline:             payload.Tagline,
		Description:         payload.Description,
		RegistrationOpensAt: payload.RegistrationOpensAt,
		StartsAt:            payload.StartsAt,
		EndsAt:              payload.EndsAt,
		LocationName:        payload.LocationName,
		LocationAddress:     payload.LocationAddress,
		LogoURL:             payload.LogoURL,
		BannerURL:           payload.BannerURL,
		ContactEmail:        payload.ContactEmail,
	}

	return c.editions.Patch(ctx, existing.ID, edition)
}
