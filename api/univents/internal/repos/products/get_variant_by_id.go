package products

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) GetVariantByID(ctx context.Context, id uuid.UUID) (*models.ProductVariant, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProductsRepo.GetVariantByID")
	defer span.End()
	result, err := database.Queries(ctx, repo.q).GetProductVariantByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapVariant(result)), nil
}
