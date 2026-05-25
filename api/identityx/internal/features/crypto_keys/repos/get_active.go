package repos

import (
	"IdentityX/internal/database/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *repo) GetActive(ctx context.Context, keyType models.CryptoKeyType, projectID *uuid.UUID) (*models.CryptoKey, error) {
	ctx, span := database.Span(ctx, repo.tracer, "GetActive")
	defer span.End()
	sqlcKeyPair, err := database.Queries(ctx, repo.q).GetActiveCryptoKey(ctx, sqlc.GetActiveCryptoKeyParams{
		Type:      string(keyType),
		ProjectID: projectID,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapKeys(sqlcKeyPair)), nil
}
