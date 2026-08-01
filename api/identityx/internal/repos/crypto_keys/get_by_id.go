package crypto_keys

import (
	"IdentityX/models"
	"context"
	"lib/database"

	"lib/telemetry"

	"github.com/google/uuid"
)

func (repo *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.CryptoKey, error) {
	ctx, span := telemetry.StartSpan(ctx, "GetByID")
	defer span.End()

	sqlcKeyPair, err := database.Queries(ctx, repo.q).GetCryptoKeyByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapKeys(sqlcKeyPair)), nil
}
