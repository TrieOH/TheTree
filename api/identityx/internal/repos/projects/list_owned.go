package projects

import (
	"IdentityX/models"
	"context"
	"lib/database"
	"lib/xslices"

	"lib/telemetry"

	"github.com/google/uuid"
)

func (repo *Repo) ListOwned(ctx context.Context, userID uuid.UUID) ([]models.Project, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListOwned")
	defer span.End()

	sqlcProjects, err := database.Queries(ctx, repo.q).ListOwnedProjects(ctx, userID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcProjects, mapProject), nil
}
