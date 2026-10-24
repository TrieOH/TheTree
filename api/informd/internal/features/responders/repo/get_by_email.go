package repo

import (
	"context"

	"Informd/models"
	"lib/database"
)

func (repo *Repo) GetByEmail(ctx context.Context, email string) (*models.Responder, error) {
	row, err := database.Queries(ctx, repo.q).GetResponderByEmail(ctx, email)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapResponder(row)), nil
}
