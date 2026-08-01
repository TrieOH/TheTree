package commands

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"
	"univents/internal/authz"
	"univents/models"
)

func (c *Commands) PatchProduct(ctx context.Context, payload models.PatchProductInput) (*models.Product, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProductsService.PatchProduct")
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

	err = authz.Service.CheckEvent(ctx, ident.Sub.ID, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	product := &models.Product{
		VendorCode:           payload.VendorCode,
		RequiresRegistration: payload.RequiresRegistration,
	}

	return c.products.PatchProduct(ctx, payload.ProductID, product)
}
