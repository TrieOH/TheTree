package repo

import (
	"IdentityX/models"
	"context"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *Repo) ListJoined(ctx context.Context, userID uuid.UUID) ([]models.Project, error) {
	sqlcProjects, err := database.Queries(ctx, repo.q).ListJoinedProjects(ctx, userID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcProjects, mapProject), nil
}
