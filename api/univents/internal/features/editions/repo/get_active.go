package repo

import (
	"context"
	"lib/database"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) GetActive(ctx context.Context, eventID uuid.UUID) (*models.Edition, error) {
	edition, err := database.Queries(ctx, repo.q).GetActiveEdition(ctx, eventID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapEdition(edition)), nil
}
