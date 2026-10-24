package service

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) DeleteVariant(ctx context.Context, id uuid.UUID) error {
	existing, err := o.Products.GetVariantByID(ctx, id)
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

	return o.Products.DeleteVariant(ctx, id)
}
