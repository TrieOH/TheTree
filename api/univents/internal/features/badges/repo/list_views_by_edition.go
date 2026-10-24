package repo

import (
	"context"
	"lib/database"
	"lib/xslices"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) ListViewsByEdition(ctx context.Context, editionID uuid.UUID) ([]models.BadgeEmissionView, error) {
	rows, err := database.Queries(ctx, repo.q).ListBadgeEmissionViewsByEdition(ctx, editionID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(rows, mapBadgeEmissionViewFromEdition), nil
}
