package repos

import (
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *editionsRepo) Finish(ctx context.Context, editionID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "EditionsRepo.Finish")
	defer span.End()

	err := database.Queries(ctx, repo.q).FinishEdition(ctx, editionID)
	if err != nil {
		return repo.dbe(err)
	}

	return nil
}
