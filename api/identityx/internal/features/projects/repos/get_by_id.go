package repos

import (
	"IdentityX/models"
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	ctx, span := repo.tracer.Start(ctx, "GetByID")
	defer span.End()
	sqlcProject, err := database.Queries(ctx, repo.q).GetProjectByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapProject(sqlcProject)), nil
}
