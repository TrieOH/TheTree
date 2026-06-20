package repos

import (
	"IdentityX/models"
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *repo) GetProjectServiceAccount(ctx context.Context, id uuid.UUID) (*models.Actor, error) {
	ctx, span := database.Span(ctx, repo.tracer, "GetProjectServiceAccount")
	defer span.End()
	sqlcActor, err := database.Queries(ctx, repo.q).GetProjectServiceAccount(ctx, &id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapActor(sqlcActor)), nil
}
