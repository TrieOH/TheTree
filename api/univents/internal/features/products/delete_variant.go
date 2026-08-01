package products

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"
	"univents/internal/authz"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) DeleteVariant(ctx context.Context, id uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "ProductsService.DeleteVariant")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	existing, err := o.products.GetVariantByID(ctx, id)
	if err != nil {
		return err
	}

	edition, err := o.editions.GetByID(ctx, existing.EditionID)
	if err != nil {
		return err
	}

	err = authz.Service.CheckEvent(ctx, ident.Sub.ID, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return err
	}

	return o.products.DeleteVariant(ctx, id)
}
