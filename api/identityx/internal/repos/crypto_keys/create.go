package crypto_keys

import (
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"context"
	"lib/crypto"
	"lib/database"
	"time"

	"lib/telemetry"

	"github.com/google/uuid"
)

func (repo *Repo) Create(ctx context.Context, projectID *uuid.UUID, pair *crypto.KeyPair, keyType string, expiresAt *time.Time) (*models.CryptoKey, error) {
	ctx, span := telemetry.StartSpan(ctx, "Create")
	defer span.End()

	sqlcKeyPair, err := database.Queries(ctx, repo.q).CreateCryptoKey(ctx, sqlc.CreateCryptoKeyParams{
		ProjectID:           projectID,
		Type:                keyType,
		PublicKey:           pair.Public,
		EncryptedPrivateKey: pair.EncryptedPrivate,
		Algorithm:           pair.Algorithm,
		ExpiresAt:           expiresAt,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapKeys(sqlcKeyPair)), nil
}
