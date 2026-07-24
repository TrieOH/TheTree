package repos

import (
	"context"
	"lib/database"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *repo) GetProductByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	ctx, span := repo.tracer.Start(ctx, "ProductsRepo.GetProductByID")
	defer span.End()
	result, err := database.Queries(ctx, repo.q).GetProductByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapProduct(result)), nil
}
