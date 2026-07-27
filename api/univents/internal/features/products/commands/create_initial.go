package commands

import (
	"context"
	"lib/database"
	"lib/telemetry"
	idx "sdk/identityx"
	"univents/models"

	"github.com/MintzyG/fun"
)

func (c *Commands) CreateInitial(ctx context.Context, payload models.CreateInitialProductInput) (*models.Product, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProductsService.CreateInitial")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	edition, err := c.editions.GetByID(ctx, payload.EditionID)
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

	var product *models.Product
	err = database.RunTx(ctx, func(ctx context.Context) error {
		p, err := c.products.CreateProduct(ctx, &models.Product{
			EditionID:            payload.EditionID,
			VendorCode:           payload.VendorCode,
			RequiresRegistration: payload.RequiresRegistration,
		})
		if err != nil {
			return err
		}

		_, err = c.products.CreateVariant(ctx, &models.ProductVariant{
			EditionID:   payload.EditionID,
			ProductID:   p.ID,
			VendorCode:  payload.VariantVendorCode,
			Name:        payload.Name,
			Description: payload.Description,
			Price:       payload.Price,
			Stock:       payload.Stock,
		})
		if err != nil {
			return err
		}

		product = p
		return nil
	})
	if err != nil {
		return nil, err
	}

	return product, nil
}
