package products

import (
	"context"
	"lib/database"
	"lib/telemetry"
	idx "sdk/identityx"
	"univents/models"
)

func (o *Operations) CreateInitial(ctx context.Context, payload models.CreateInitialProductInput) (*models.Product, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProductsService.CreateInitial")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	edition, err := o.editions.GetByID(ctx, payload.EditionID)
	if err != nil {
		return nil, err
	}

	err = o.authz.CheckEvent(ctx, ident.Sub.ID, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	var product *models.Product
	err = database.RunTx(ctx, func(ctx context.Context) error {
		p, err := o.products.CreateProduct(ctx, &models.Product{
			EditionID:            payload.EditionID,
			VendorCode:           payload.VendorCode,
			RequiresRegistration: payload.RequiresRegistration,
		})
		if err != nil {
			return err
		}

		_, err = o.products.CreateVariant(ctx, &models.ProductVariant{
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
