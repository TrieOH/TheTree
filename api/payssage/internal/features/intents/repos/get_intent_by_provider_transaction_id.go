package repos

import (
	"context"
	"lib/database"
	"payssage/internal/sqlc"
	"payssage/models"
)

func (repo *Repo) GetByProviderTransactionID(ctx context.Context, provider string, transactionID string) (*models.Intent, error) {
	ctx, span := repo.tracer.Start(ctx, "IntentRepo.GetByProviderTransactionID")
	defer span.End()

	sqlcIntent, err := database.Queries(ctx, repo.q).GetIntentByProviderTransactionID(ctx, sqlc.GetIntentByProviderTransactionIDParams{
		Provider:      provider,
		TransactionID: transactionID,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapIntent(sqlcIntent)), nil
}
