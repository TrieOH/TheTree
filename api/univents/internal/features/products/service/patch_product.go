package service

import (
	"context"
	"univents/models"
)

func (o *Operations) PatchProduct(ctx context.Context, payload models.PatchProductInput) (*models.Product, error) {
	existing, err := o.Products.GetProductByID(ctx, payload.ProductID)
	if err != nil {
		return nil, err
	}

	edition, err := o.Editions.GetByID(ctx, existing.EditionID)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Event(ctx, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	product := &models.Product{
		VendorCode:           payload.VendorCode,
		RequiresRegistration: payload.RequiresRegistration,
	}

	return o.Products.PatchProduct(ctx, payload.ProductID, product)
}
