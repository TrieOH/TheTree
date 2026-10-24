package service

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) ListProductsByEdition(ctx context.Context, editionID uuid.UUID) ([]models.Product, error) {
	return o.Products.ListProductsByEdition(ctx, editionID)
}
