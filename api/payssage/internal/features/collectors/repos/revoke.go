package repos

import (
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *repo) Revoke(ctx context.Context, id uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "CollectorRepo.Revoke")
	defer span.End()

	return repo.dbe(database.Queries(ctx, repo.q).RevokeCollector(ctx, id))
}
