package repos

import (
	"IdentityX/models"
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Actor, error) {
	ctx, span := database.Span(ctx, repo.tracer, "GetByID")
	defer span.End()
	sqlcActor, err := database.Queries(ctx, repo.q).GetActorByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapActor(sqlcActor)), nil
}
