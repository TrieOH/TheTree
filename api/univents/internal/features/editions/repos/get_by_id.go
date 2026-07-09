package repos

import (
	"context"
	"lib/database"
	"univents/contracts"

	"github.com/google/uuid"
)

func (repo *editionsRepo) GetByID(ctx context.Context, editionID uuid.UUID) (*contracts.Edition, error) {
	ctx, span := repo.tracer.Start(ctx, "EditionsRepo.GetByID")
	defer span.End()

	sqlcEdition, err := database.Queries(ctx, repo.q).GetEditionByID(ctx, editionID)
	if err != nil {
		return nil, repo.dbe(err)
	}

	return mapEditionFromDB(&sqlcEdition), nil
}
