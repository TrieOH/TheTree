package commands

import (
	"context"
	idx "sdk/identityx"
	"univents/models"

	"github.com/MintzyG/fun"
)

func (c *Commands) PatchProduct(ctx context.Context, payload models.PatchProductInput) (*models.Product, error) {
	ctx, span := c.tracer.Start(ctx, "ProductsService.PatchProduct")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	existing, err := c.products.GetProductByID(ctx, payload.ProductID)
	if err != nil {
		return nil, err
	}

	edition, err := c.editions.GetByID(ctx, existing.EditionID)
	if err != nil {
		return nil, err
	}

	member, err := c.events.GetMember(ctx, edition.EventID, ident.Sub.ID)
	if fun.Is(err, fun.CodeNotFound) {
		return nil, fun.ErrForbidden("insufficient permissions")
	}
	if err != nil {
		return nil, err
	}
	if !member.Role.Minimum(models.EventMemberRoleAdmin) {
		return nil, fun.ErrForbidden("insufficient permissions")
	}

	product := &models.Product{
		VendorCode:           payload.VendorCode,
		RequiresRegistration: payload.RequiresRegistration,
	}

	return c.products.PatchProduct(ctx, payload.ProductID, product)
}
