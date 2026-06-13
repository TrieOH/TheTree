package repos

import (
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *repo) UpdateLastLoginAt(ctx context.Context, actorID uuid.UUID) error {
	ctx, span := database.Span(ctx, repo.tracer, "UpdateLastLoginAt")
	defer span.End()
	return repo.dbe(database.Queries(ctx, repo.q).UpdateActorLastLoginAt(ctx, actorID))
}
