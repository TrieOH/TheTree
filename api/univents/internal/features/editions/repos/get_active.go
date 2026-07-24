package repos

import (
	"context"
	"lib/database"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *repo) GetActive(ctx context.Context, eventID uuid.UUID) (*models.Edition, error) {
	ctx, span := repo.tracer.Start(ctx, "EditionsRepo.GetActive")
	defer span.End()
	edition, err := database.Queries(ctx, repo.q).GetActiveEdition(ctx, eventID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapEdition(edition)), nil
}
