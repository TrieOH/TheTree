package repos

import (
	"context"
	"lib/database"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (repo *Repo) Revoke(ctx context.Context, id uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "CollectorRepo.Revoke")
	defer span.End()

	return repo.dbe(database.Queries(ctx, repo.q).RevokeCollector(ctx, id))
}
