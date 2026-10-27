package repo

import (
	"context"
	"lib/database"
	"payssage/models"

	"github.com/google/uuid"
)

func (repo *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Intent, error) {
	sqlcIntent, err := database.Queries(ctx, repo.q).GetIntentByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}

	return new(mapIntent(sqlcIntent)), nil
}
