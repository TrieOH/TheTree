package repos

import (
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *editionsRepo) Open(ctx context.Context, editionID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "EditionsRepo.Open")
	defer span.End()

	err := database.Queries(ctx, repo.q).OpenEditionRegistrations(ctx, editionID)
	if err != nil {
		return repo.dbe(err)
	}

	return nil
}
