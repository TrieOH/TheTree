package repos

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"payssage/internal/sqlc"
	"payssage/models"
)

func (repo *Repo) Update(ctx context.Context, toUpdate models.Intent) (*models.Intent, error) {
	ctx, span := telemetry.StartSpan(ctx, "IntentRepo.Update")
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
