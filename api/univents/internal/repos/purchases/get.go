package purchases

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"univents/internal/sqlc"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Purchase, error) {
	ctx, span := telemetry.StartSpan(ctx, "PurchasesRepo.GetByID")
	defer span.End()
	result, err := database.Queries(ctx, repo.q).GetPurchaseByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapPurchase(result)), nil
}

func (repo *Repo) GetByIDForOwner(ctx context.Context, id, purchaserID uuid.UUID) (*models.Purchase, error) {
	ctx, span := telemetry.StartSpan(ctx, "PurchasesRepo.GetByIDForOwner")
	defer span.End()
	result, err := database.Queries(ctx, repo.q).GetPurchaseByIDForOwner(ctx, sqlc.GetPurchaseByIDForOwnerParams{
		ID:          id,
		PurchaserID: purchaserID,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapPurchase(result)), nil
}

func (repo *Repo) GetByIntentID(ctx context.Context, intentID uuid.UUID) (*models.Purchase, error) {
	ctx, span := telemetry.StartSpan(ctx, "PurchasesRepo.GetByIntentID")
	defer span.End()
	result, err := database.Queries(ctx, repo.q).GetPurchaseByIntentID(ctx, &intentID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapPurchase(result)), nil
}

func (repo *Repo) ListByPurchaser(ctx context.Context, purchaserID uuid.UUID) ([]models.Purchase, error) {
	ctx, span := telemetry.StartSpan(ctx, "PurchasesRepo.ListByPurchaser")
	defer span.End()
	rows, err := database.Queries(ctx, repo.q).ListPurchasesByPurchaser(ctx, purchaserID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	out := make([]models.Purchase, 0, len(rows))
	for _, row := range rows {
		out = append(out, mapPurchase(row))
	}
	return out, nil
}

func (repo *Repo) ListItemsByPurchase(ctx context.Context, purchaseID uuid.UUID) ([]models.PurchaseItem, error) {
	ctx, span := telemetry.StartSpan(ctx, "PurchasesRepo.ListItemsByPurchase")
	defer span.End()
	rows, err := database.Queries(ctx, repo.q).ListPurchaseItemsByPurchase(ctx, purchaseID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	out := make([]models.PurchaseItem, 0, len(rows))
	for _, row := range rows {
		out = append(out, mapPurchaseItem(row))
	}
	return out, nil
}
