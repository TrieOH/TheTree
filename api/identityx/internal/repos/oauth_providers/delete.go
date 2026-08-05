package oauth_providers

import (
	"context"
	"lib/database"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (repo *Repo) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "Delete")
	defer span.End()

	err := database.Queries(ctx, repo.q).DeleteProjectOAuthProvider(ctx, id)
	if err != nil {
		return repo.dbe(err)
	}
	return nil
}
