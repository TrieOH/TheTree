package queries

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

func (q *Queries) ListProductsByEdition(ctx context.Context, editionID uuid.UUID) ([]models.Product, error) {
	ctx, span := q.tracer.Start(ctx, "ProductsService.ListProductsByEdition")
	defer span.End()
	return q.products.ListProductsByEdition(ctx, editionID)
}
