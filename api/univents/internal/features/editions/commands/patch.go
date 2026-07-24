package commands

import (
	"context"
	idx "sdk/identityx"
	"univents/models"

	"github.com/MintzyG/fun"
)

func (c *Commands) Patch(ctx context.Context, payload models.PatchEditionInput) (*models.Edition, error) {
	ctx, span := c.tracer.Start(ctx, "EditionService.Patch")
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

	member, err := c.events.GetMember(ctx, event.ID, ident.Sub.ID)
	if fun.Is(err, fun.CodeNotFound) {
		return nil, fun.ErrForbidden("insufficient permissions")
	}
	if err != nil {
		return nil, err
	}
	if !member.Role.Minimum(models.EventMemberRoleAdmin) {
		return nil, fun.ErrForbidden("insufficient permissions")
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
