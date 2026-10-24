package service

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "ProductsService.DeleteProduct")
	defer span.End()

	existing, err := o.Products.GetProductByID(ctx, id)
	if err != nil {
		return err
	}

	edition, err := o.Editions.GetByID(ctx, existing.EditionID)
	if err != nil {
		return err
	}

	_, err = o.Authz.Event(ctx, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return err
	}

	return o.Products.DeleteProduct(ctx, id)
}
