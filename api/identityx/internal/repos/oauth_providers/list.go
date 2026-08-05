package oauth_providers

import (
	"IdentityX/models"
	"context"
	"lib/database"
	"lib/telemetry"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *Repo) ListByProject(ctx context.Context, projectID uuid.UUID) ([]models.ProjectOAuthProviders, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListByProject")
	defer span.End()

	rows, err := database.Queries(ctx, repo.q).ListProjectOAuthProviders(ctx, projectID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(rows, mapProjectOAuthProvider), nil
}
