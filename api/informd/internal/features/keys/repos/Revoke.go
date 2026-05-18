package repos

import (
	"Informd/internal/database/sqlc"
	"Informd/models"
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *repo) Revoke(ctx context.Context, id, userID uuid.UUID) (*models.APIKey, error) {
	ctx, span := repo.tracer.Start(ctx, "Revoke")
	defer span.End()
	sqlcApiKey, err := database.Queries(ctx, repo.q).RevokeAPIKey(ctx, sqlc.RevokeAPIKeyParams{
		ID:      id,
		OwnerID: userID,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapApiKey(sqlcApiKey)), nil
}
