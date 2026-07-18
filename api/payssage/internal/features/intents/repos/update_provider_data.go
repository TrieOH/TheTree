package repos

import (
	"context"
	"encoding/json"
	"lib/database"
	"payssage/internal/database/sqlc"
	"payssage/models"

	"github.com/google/uuid"
)

func (repo *repo) UpdateProviderData(ctx context.Context, id uuid.UUID, providerData json.RawMessage) (*models.Intent, error) {
	ctx, span := repo.tracer.Start(ctx, "IntentRepo.UpdateProviderData")
	defer span.End()

	sqlcIntent, err := database.Queries(ctx, repo.q).UpdateIntentProviderData(ctx, sqlc.UpdateIntentProviderDataParams{
		ID:           id,
		ProviderData: providerData,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}

	return new(mapIntent(sqlcIntent)), nil
}
