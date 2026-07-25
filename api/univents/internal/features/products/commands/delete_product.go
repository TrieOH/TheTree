package commands

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"
	"univents/models"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (c *Commands) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "ProductsService.DeleteProduct")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	existing, err := c.products.GetProductByID(ctx, id)
	if err != nil {
		return err
	}

	edition, err := c.editions.GetByID(ctx, existing.EditionID)
	if err != nil {
		return err
	}

	member, err := c.events.GetMember(ctx, edition.EventID, ident.Sub.ID)
	if fun.Is(err, fun.CodeNotFound) {
		return fun.ErrForbidden("insufficient permissions")
	}
	if err != nil {
		return err
	}
	if !member.Role.Minimum(models.EventMemberRoleAdmin) {
		return fun.ErrForbidden("insufficient permissions")
	}

	return c.products.DeleteProduct(ctx, id)
}
