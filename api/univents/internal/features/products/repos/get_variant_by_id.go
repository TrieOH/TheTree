package repos

import (
	"context"
	"lib/database"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *repo) GetVariantByID(ctx context.Context, id uuid.UUID) (*models.ProductVariant, error) {
	ctx, span := repo.tracer.Start(ctx, "ProductsRepo.GetVariantByID")
	defer span.End()
	result, err := database.Queries(ctx, repo.q).GetProductVariantByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapVariant(result)), nil
}
