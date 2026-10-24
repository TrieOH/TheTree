package repo

import (
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *Repo) UpdateLastLoginAt(ctx context.Context, actorID uuid.UUID) error {
	return repo.dbe(database.Queries(ctx, repo.q).UpdateActorLastLoginAt(ctx, actorID))
}
