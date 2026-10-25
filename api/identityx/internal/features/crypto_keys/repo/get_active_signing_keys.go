package repo

import (
	"IdentityX/models"
	"context"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *Repo) GetActiveSigningKeys(ctx context.Context, projectID *uuid.UUID) ([]models.ActiveSigningKey, error) {
	sqlcKeys, err := database.Queries(ctx, repo.q).GetActiveSigningKeys(ctx, projectID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcKeys, mapToActiveSigningKey), nil
}
