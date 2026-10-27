package repo

import (
	"IdentityX/models"
	"context"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *Repo) List(ctx context.Context, projectID uuid.UUID) ([]models.Actor, error) {
	sqlcActors, err := database.Queries(ctx, repo.q).ListActorsFromProject(ctx, &projectID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcActors, mapActor), nil
}
