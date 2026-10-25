package repo

import (
	"context"

	"Informd/models"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *Repo) GetSelectConfig(ctx context.Context, fieldID uuid.UUID) (*models.FieldSelectConfig, error) {
	sqlcConfig, err := database.Queries(ctx, repo.q).GetFieldSelectConfig(ctx, fieldID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapFieldSelectConfig(sqlcConfig)), nil
}
