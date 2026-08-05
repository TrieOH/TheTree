package oauth_providers

import (
	"context"
	"lib/database"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (repo *Repo) DeleteState(ctx context.Context, id uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "DeleteState")
	defer span.End()

	err := database.Queries(ctx, repo.q).DeleteOAuthLoginState(ctx, id)
	if err != nil {
		return repo.dbe(err)
	}
	return nil
}
