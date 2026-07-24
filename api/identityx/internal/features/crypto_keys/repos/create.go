package repos

import (
	"IdentityX/internal/database/sqlc"
	"IdentityX/models"
	"context"
	"lib/crypto"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *Repo) Create(ctx context.Context, projectID *uuid.UUID, pair *crypto.KeyPair, keyType string) (*models.CryptoKey, error) {
	ctx, span := database.Span(ctx, repo.tracer, "Create")
	defer span.End()
	sqlcKeyPair, err := database.Queries(ctx, repo.q).CreateCryptoKey(ctx, sqlc.CreateCryptoKeyParams{
		ProjectID:           projectID,
		Type:                keyType,
		PublicKey:           pair.Public,
		EncryptedPrivateKey: pair.EncryptedPrivate,
		Algorithm:           pair.Algorithm,
		ExpiresAt:           nil,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapKeys(sqlcKeyPair)), nil
}
