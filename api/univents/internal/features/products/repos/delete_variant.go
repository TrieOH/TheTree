package repos

import (
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *repo) DeleteVariant(ctx context.Context, id uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "ProductsRepo.DeleteVariant")
	defer span.End()
	err := database.Queries(ctx, repo.q).DeleteProductVariant(ctx, id)
	if err != nil {
		return repo.dbe(err)
	}
	return nil
}
