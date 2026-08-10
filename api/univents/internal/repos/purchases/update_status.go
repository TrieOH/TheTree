package purchases

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"univents/internal/sqlc"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) UpdateStatus(ctx context.Context, id uuid.UUID, status models.PurchaseStatus, reason *string) (*models.Purchase, error) {
	ctx, span := telemetry.StartSpan(ctx, "PurchasesRepo.UpdateStatus")
	defer span.End()
	result, err := database.Queries(ctx, repo.q).UpdatePurchaseStatus(ctx, sqlc.UpdatePurchaseStatusParams{
		ID:           id,
		Status:       string(status),
		StatusReason: reason,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapPurchase(result)), nil
}
