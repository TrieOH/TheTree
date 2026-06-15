package repos

import (
	"IdentityX/models"
	"context"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *repo) ListOwned(ctx context.Context, userID uuid.UUID) ([]models.Project, error) {
	ctx, span := repo.tracer.Start(ctx, "ListOwned")
	defer span.End()
	sqlcProjects, err := database.Queries(ctx, repo.q).ListOwnedProjects(ctx, userID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcProjects, mapProject), nil
}
