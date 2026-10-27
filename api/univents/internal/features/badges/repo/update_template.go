package repo

import (
	"context"
	"encoding/json"
	"lib/database"
	"univents/internal/sqlc"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) Update(ctx context.Context, id uuid.UUID, name string, designData json.RawMessage) (*models.BadgeTemplate, error) {
	row, err := database.Queries(ctx, repo.q).UpdateBadgeTemplate(ctx, sqlc.UpdateBadgeTemplateParams{
		ID:         id,
		Name:       name,
		DesignData: designData,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapBadgeTemplate(row)), nil
}
