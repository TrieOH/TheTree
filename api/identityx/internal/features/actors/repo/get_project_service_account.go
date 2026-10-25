package repo

import (
	"IdentityX/models"
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *Repo) GetProjectServiceAccount(ctx context.Context, id uuid.UUID) (*models.Actor, error) {
	sqlcActor, err := database.Queries(ctx, repo.q).GetProjectServiceAccount(ctx, &id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapActor(sqlcActor)), nil
}
