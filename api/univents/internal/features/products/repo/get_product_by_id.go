package repo

import (
	"context"
	"lib/database"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) GetProductByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	result, err := database.Queries(ctx, repo.q).GetProductByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapProduct(result)), nil
}
