package repos

import (
	"context"
	"lib/database"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (repo *Repo) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "ProductsRepo.DeleteProduct")
	defer span.End()
	err := database.Queries(ctx, repo.q).DeleteProduct(ctx, id)
	if err != nil {
		return repo.dbe(err)
	}
	return nil
}
