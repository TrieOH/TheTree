package repos

import (
	"IdentityX/models"
	"context"
	"lib/database"
	"lib/xslices"

	"lib/telemetry"

	"github.com/google/uuid"
)

func (repo *Repo) List(ctx context.Context, projectID uuid.UUID) ([]models.Actor, error) {
	ctx, span := telemetry.StartSpan(ctx, "List")
	defer span.End()

	sqlcActors, err := database.Queries(ctx, repo.q).ListActorsFromProject(ctx, &projectID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcActors, mapActor), nil
}
