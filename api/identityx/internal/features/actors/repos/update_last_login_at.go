package repos

import (
	"context"
	"lib/database"

	"lib/telemetry"

	"github.com/google/uuid"
)

func (repo *Repo) UpdateLastLoginAt(ctx context.Context, actorID uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "UpdateLastLoginAt")
	defer span.End()

	return repo.dbe(database.Queries(ctx, repo.q).UpdateActorLastLoginAt(ctx, actorID))
}
