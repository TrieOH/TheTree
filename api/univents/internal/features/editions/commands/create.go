package commands

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"
	"univents/models"

	"github.com/MintzyG/fun"
)

func (c *Commands) Create(ctx context.Context, payload models.CreateEditionInput) (*models.Edition, error) {
	ctx, span := telemetry.StartSpan(ctx, "EditionService.Create")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	event, err := c.events.GetByID(ctx, payload.EventID)
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
		EventID:   payload.EventID,
		Name:      payload.Name,
		Slug:      payload.Slug,
		IsDraft:   true,
		StartsAt:  payload.StartsAt,
		EndsAt:    payload.EndsAt,
		CreatedBy: ident.Sub.ID,
	}

	return c.editions.Create(ctx, edition)
}
