package service

import (
	"context"
	"univents/models"
)

func (o *Operations) CreateVariant(ctx context.Context, payload models.CreateProductVariantInput) (*models.ProductVariant, error) {
	product, err := o.Products.GetProductByID(ctx, payload.ProductID)
	if err != nil {
		return nil, err
	}

	edition, err := o.Editions.GetByID(ctx, product.EditionID)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Event(ctx, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	variant := &models.ProductVariant{
		EditionID:   product.EditionID,
		ProductID:   product.ID,
		VendorCode:  payload.VendorCode,
		Name:        payload.Name,
		Description: payload.Description,
		Price:       payload.Price,
		Stock:       payload.Stock,
	}

	return o.Products.CreateVariant(ctx, variant)
}
