package repos

import (
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *editionsRepo) Announce(ctx context.Context, editionID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "EditionsRepo.Announce")
	defer span.End()

	err := database.Queries(ctx, repo.q).AnnounceEdition(ctx, editionID)
	if err != nil {
		return repo.dbe(err)
	}

	return nil
}
