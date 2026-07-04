package repos

import (
	"IdentityX/models"
	"context"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *repo) List(ctx context.Context, projectID uuid.UUID) ([]models.Actor, error) {
	ctx, span := database.Span(ctx, repo.tracer, "List")
	defer span.End()
	sqlcActors, err := database.Queries(ctx, repo.q).ListActorsFromProject(ctx, &projectID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcActors, mapActor), nil
}
