package repos

import (
	"context"
	"lib/database"
	"payssage/internal/database/sqlc"
	"payssage/models"
)

func (repo *repo) Update(ctx context.Context, toUpdate models.Intent) (*models.Intent, error) {
	ctx, span := repo.tracer.Start(ctx, "IntentRepo.Update")
	defer span.End()

	var statusDetail *string
	if toUpdate.StatusDetail != nil {
		statusDetail = new(string(*toUpdate.StatusDetail))
	}

	sqlcIntent, err := database.Queries(ctx, repo.q).UpdateIntent(ctx, sqlc.UpdateIntentParams{
		ID:           toUpdate.ID,
		Status:       string(toUpdate.Status),
		StatusDetail: statusDetail,
		ProviderData: toUpdate.ProviderData,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}

	return new(mapIntent(sqlcIntent)), nil
}
