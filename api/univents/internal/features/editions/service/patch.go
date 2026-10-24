package service

import (
	"context"
	"univents/models"
)

func (o *Operations) Patch(ctx context.Context, payload models.PatchEditionInput) (*models.Edition, error) {
	existing, err := o.Editions.GetByID(ctx, payload.EditionID)
	if err != nil {
		return nil, err
	}

	event, err := o.Events.GetByID(ctx, existing.EventID)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Event(ctx, event.ID, models.EventMemberRoleAdmin)
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

	return o.Editions.Patch(ctx, existing.ID, edition)
}
