package crypto_keys

import (
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"

	"lib/telemetry"

	"github.com/google/uuid"
)

func (repo *Repo) GetActive(ctx context.Context, keyType models.CryptoKeyType, projectID *uuid.UUID) (*models.CryptoKey, error) {
	ctx, span := telemetry.StartSpan(ctx, "GetActive")
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
