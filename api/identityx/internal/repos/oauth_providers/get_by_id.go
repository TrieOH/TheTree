package oauth_providers

import (
	"IdentityX/models"
	"context"
	"lib/database"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (repo *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.ProjectOAuthProviders, error) {
	ctx, span := telemetry.StartSpan(ctx, "GetByID")
	defer span.End()

	row, err := database.Queries(ctx, repo.q).GetProjectOAuthProvider(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapProjectOAuthProvider(row)), nil
}
