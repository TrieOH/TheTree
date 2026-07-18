package repos

import (
	"context"
	"lib/database"
	"payssage/models"

	"github.com/google/uuid"
)

func (repo *repo) Confirm(ctx context.Context, id uuid.UUID) (*models.Intent, error) {
	ctx, span := repo.tracer.Start(ctx, "IntentRepo.Confirm")
	defer span.End()

	sqlcIntent, err := database.Queries(ctx, repo.q).ConfirmIntent(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}

	return new(mapIntent(sqlcIntent)), nil
}
