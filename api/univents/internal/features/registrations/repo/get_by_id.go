package repo

import (
	"context"
	"lib/database"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Registration, error) {
	row, err := database.Queries(ctx, repo.q).GetRegistrationByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapRegistration(row)), nil
}
