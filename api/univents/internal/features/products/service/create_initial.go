package service

import (
	"context"
	"lib/telemetry"
	"univents/models"
)

func (o *Operations) CreateInitial(ctx context.Context, payload models.CreateInitialProductInput) (*models.Product, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProductsService.CreateInitial")
	defer span.End()

	edition, err := o.Editions.GetByID(ctx, payload.EditionID)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Event(ctx, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	var product *models.Product
	err = o.Tx.WithinTx(ctx, func(ctx context.Context) error {
		p, err := o.Products.CreateProduct(ctx, &models.Product{
			EditionID:            payload.EditionID,
			VendorCode:           payload.VendorCode,
			RequiresRegistration: payload.RequiresRegistration,
		})
		if err != nil {
			return err
		}

		_, err = o.Products.CreateVariant(ctx, &models.ProductVariant{
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
