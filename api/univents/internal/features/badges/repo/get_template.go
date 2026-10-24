package repo

import (
	"context"
	"lib/database"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.BadgeTemplate, error) {
	row, err := database.Queries(ctx, repo.q).GetBadgeTemplateByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}

	return new(mapBadgeTemplate(row)), nil
}
