package repos

import (
	"IdentityX/models"
	"context"
	"lib/database"

	"lib/telemetry"

	"github.com/google/uuid"
)

func (repo *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	ctx, span := telemetry.StartSpan(ctx, "GetByID")
	defer span.End()

	sqlcProject, err := database.Queries(ctx, repo.q).GetProjectByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapProject(sqlcProject)), nil
}
