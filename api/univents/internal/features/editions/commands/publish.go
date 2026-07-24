package commands

import (
	"context"
	idx "sdk/identityx"
	"univents/models"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (c *Commands) Publish(ctx context.Context, editionID uuid.UUID) error {
	ctx, span := c.tracer.Start(ctx, "EditionService.Publish")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	edition, err := c.editions.GetByID(ctx, editionID)
	if err != nil {
		return err
	}

	if !edition.IsDraft {
		return fun.ErrBadRequest("edition is already published")
	}

	event, err := c.events.GetByID(ctx, edition.EventID)
	if err != nil {
		return err
	}

	member, err := c.events.GetMember(ctx, event.ID, ident.Sub.ID)
	if fun.Is(err, fun.CodeNotFound) {
		return fun.ErrForbidden("insufficient permissions")
	}
	if err != nil {
		return err
	}
	if !member.Role.Minimum(models.EventMemberRoleAdmin) {
		return fun.ErrForbidden("insufficient permissions")
	}

	return c.editions.Publish(ctx, editionID)
}
