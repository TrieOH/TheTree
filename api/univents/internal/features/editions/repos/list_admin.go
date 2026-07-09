package repos

import (
	"context"
	"lib/database"
	"univents/contracts"

	"github.com/google/uuid"
)

func (repo *editionsRepo) ListAdmin(ctx context.Context, eventID uuid.UUID) ([]contracts.Edition, error) {
	ctx, span := repo.tracer.Start(ctx, "EditionsRepo.ListAdmin")
	defer span.End()

	sqlcEditions, err := database.Queries(ctx, repo.q).ListEditionsAdmin(ctx, eventID)
	if err != nil {
		return nil, repo.dbe(err)
	}

	outEditions := make([]contracts.Edition, 0, len(sqlcEditions))
	for _, sqlcEdition := range sqlcEditions {
		outEditions = append(outEditions, *mapEditionFromDB(&sqlcEdition))
	}
	return outEditions, nil
}
