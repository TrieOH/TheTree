package repos

import (
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *Repo) Revoke(ctx context.Context, id uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "SellerRepo.Revoke")
	defer span.End()

	return repo.dbe(database.Queries(ctx, repo.q).RevokeSeller(ctx, id))
}
