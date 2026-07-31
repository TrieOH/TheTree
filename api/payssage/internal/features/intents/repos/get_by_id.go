package repos

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"payssage/models"

	"github.com/google/uuid"
)

func (repo *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Intent, error) {
	ctx, span := telemetry.StartSpan(ctx, "IntentRepo.GetByID")
	defer span.End()

	sqlcIntent, err := database.Queries(ctx, repo.q).GetIntentByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}

	return new(mapIntent(sqlcIntent)), nil
}
