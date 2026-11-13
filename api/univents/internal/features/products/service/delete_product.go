package service

import (
	"context"
	"univents/models"
	"uuid"
)

// DeleteProduct deletes the product and returns it as it was.
func (o *Operations) DeleteProduct(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	existing, err := o.Products.GetProductByID(ctx, id)
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

	err = o.Products.DeleteProduct(ctx, id)
	if err != nil {
		return nil, err
	}
	return existing, nil
}
