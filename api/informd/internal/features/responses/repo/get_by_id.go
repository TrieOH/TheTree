package repo

import (
	"context"

	"Informd/models"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Response, error) {
	row, err := database.Queries(ctx, repo.q).GetResponseByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapResponse(row)), nil
}
