package repos

import (
	"context"
	"lib/database"
	"lib/xslices"
	"payssage/models"

	"github.com/google/uuid"
)

func (repo *Repo) ListByOwner(ctx context.Context, ownerID uuid.UUID) ([]models.Intent, error) {
	ctx, span := repo.tracer.Start(ctx, "IntentRepo.ListByOwner")
	defer span.End()

	sqlcIntents, err := database.Queries(ctx, repo.q).ListIntentsByOwner(ctx, ownerID)
	if err != nil {
		return nil, repo.dbe(err)
	}

	return xslices.MapSlice(sqlcIntents, mapIntent), nil
}
