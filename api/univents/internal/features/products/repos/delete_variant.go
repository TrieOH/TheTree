package repos

import (
	"context"
	"lib/database"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (repo *Repo) DeleteVariant(ctx context.Context, id uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "ProductsRepo.DeleteVariant")
	defer span.End()
	err := database.Queries(ctx, repo.q).DeleteProductVariant(ctx, id)
	if err != nil {
		return repo.dbe(err)
	}
	return nil
}
