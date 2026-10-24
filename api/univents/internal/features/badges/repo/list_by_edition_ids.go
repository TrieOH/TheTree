package repo

import (
	"context"
	"lib/database"
	"lib/xslices"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) ListByEditionIDs(ctx context.Context, editionIDs []uuid.UUID) ([]models.BadgeTemplate, error) {
	rows, err := database.Queries(ctx, repo.q).ListBadgeTemplatesByEditionIDs(ctx, editionIDs)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(rows, mapBadgeTemplate), nil
}
