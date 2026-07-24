package repos

import (
	"context"
	"lib/database"
	"lib/xslices"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *repo) ListProductsByEdition(ctx context.Context, editionID uuid.UUID) ([]models.Product, error) {
	ctx, span := repo.tracer.Start(ctx, "ProductsRepo.ListProductsByEdition")
	defer span.End()
	results, err := database.Queries(ctx, repo.q).ListProductsByEdition(ctx, editionID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(results, mapProduct), nil
}
