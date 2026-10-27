package repo

import (
	"context"
	"lib/database"
	"univents/models"

	"github.com/google/uuid"
)

// ListByEdition returns every purchase of an edition, newest first — the
// organizer orders read (refund plan B3).
func (repo *Repo) ListByEdition(ctx context.Context, editionID uuid.UUID) ([]models.Purchase, error) {
	rows, err := database.Queries(ctx, repo.q).ListPurchasesByEdition(ctx, editionID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	out := make([]models.Purchase, 0, len(rows))
	for _, row := range rows {
		out = append(out, mapPurchase(row))
	}
	return out, nil
}
