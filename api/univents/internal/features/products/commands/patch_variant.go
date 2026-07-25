package commands

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"
	"univents/models"

	"github.com/MintzyG/fun"
)

func (c *Commands) PatchVariant(ctx context.Context, payload models.PatchProductVariantInput) (*models.ProductVariant, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProductsService.PatchVariant")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	existing, err := c.products.GetVariantByID(ctx, payload.VariantID)
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

	variant := &models.ProductVariant{
		VendorCode:  payload.VendorCode,
		Name:        payload.Name,
		Description: payload.Description,
		Price:       payload.Price,
		Stock:       payload.Stock,
	}

	return c.products.PatchVariant(ctx, existing.ID, variant)
}
