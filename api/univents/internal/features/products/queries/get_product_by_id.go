package queries

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (q *Queries) GetProductByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProductsService.GetProductByID")
	defer span.End()
	return q.products.GetProductByID(ctx, id)
}
