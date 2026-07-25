package repos

import (
	"context"
	"lib/database"

	"github.com/google/uuid"
	"lib/telemetry"
)

func (repo *Repo) UpdateLastLoginAt(ctx context.Context, actorID uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "UpdateLastLoginAt")
	defer span.End()

	return repo.dbe(database.Queries(ctx, repo.q).UpdateActorLastLoginAt(ctx, actorID))
}
