package repos

import (
	"context"
	"lib/database"
	"lib/xslices"
	"payssage/models"

	"github.com/google/uuid"
)

func (repo *Repo) ListByWallet(ctx context.Context, walletID uuid.UUID) ([]models.Intent, error) {
	ctx, span := repo.tracer.Start(ctx, "IntentRepo.ListByWallet")
	defer span.End()

	sqlcIntents, err := database.Queries(ctx, repo.q).ListIntentsByWallet(ctx, walletID)
	if err != nil {
		return nil, repo.dbe(err)
	}

	return xslices.MapSlice(sqlcIntents, mapIntent), nil
}
