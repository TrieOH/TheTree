package service

import (
	"context"
	"univents/models"
	"uuid"
)

func (o *Operations) ListProductsByEdition(ctx context.Context, editionID uuid.UUID) ([]models.Product, error) {
	return o.Products.ListProductsByEdition(ctx, editionID)
}
