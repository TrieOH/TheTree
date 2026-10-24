package repo

import (
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *Repo) GetActive(ctx context.Context, keyType models.CryptoKeyType, projectID *uuid.UUID) (*models.CryptoKey, error) {
	sqlcKeyPair, err := database.Queries(ctx, repo.q).GetActiveCryptoKey(ctx, sqlc.GetActiveCryptoKeyParams{
		Type:      string(keyType),
		ProjectID: projectID,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapKeys(sqlcKeyPair)), nil
}
