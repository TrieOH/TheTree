package intents

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"payssage/internal/sqlc"
	"payssage/models"
)

func (repo *Repo) Create(ctx context.Context, toCreate models.Intent) (*models.Intent, error) {
	ctx, span := telemetry.StartSpan(ctx, "IntentRepo.Create")
	defer span.End()

	var statusDetail *string
	if toCreate.StatusDetail != nil {
		statusDetail = new(string(*toCreate.StatusDetail))
	}

	sqlcIntent, err := database.Queries(ctx, repo.q).CreateIntent(ctx, sqlc.CreateIntentParams{
		ID:           toCreate.ID,
		WalletID:     toCreate.WalletID,
		SellerID:     toCreate.SellerID,
		CollectorID:  toCreate.CollectorID,
		AmountCents:  toCreate.AmountCents,
		Currency:     toCreate.Currency,
		Sandbox:      toCreate.Sandbox,
		Provider:     toCreate.Provider,
		Status:       string(toCreate.Status),
		StatusDetail: statusDetail,
		ProviderData: toCreate.ProviderData,
		Metadata:     toCreate.Metadata,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}

	return new(mapIntent(sqlcIntent)), nil
}
