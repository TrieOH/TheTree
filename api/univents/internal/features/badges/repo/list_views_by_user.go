package repo

import (
	"context"
	"lib/database"
	"lib/xslices"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) ListViewsByUser(ctx context.Context, userID uuid.UUID) ([]models.BadgeEmissionView, error) {
	rows, err := database.Queries(ctx, repo.q).ListBadgeEmissionViewsByUser(ctx, userID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(rows, mapBadgeEmissionView), nil
}
