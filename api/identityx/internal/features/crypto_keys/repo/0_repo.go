package repo

import (
	"IdentityX/errs"
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"IdentityX/ports"
	"lib/database"
)

type Repo struct {
	q   *sqlc.Queries
	dbe database.NamedErrorHandler
}

var _ ports.CryptoKeysRepo = (*Repo)(nil)

func New(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewNamedErrorHandler(errs.ErrCryptoKeyNotFound),
	}
}

func mapKeys(src sqlc.CryptoKey) models.CryptoKey {
	return models.CryptoKey{
		ID:                  src.ID,
		ProjectID:           src.ProjectID,
		Type:                models.CryptoKeyType(src.Type),
		Status:              models.CryptoKeyStatus(src.Status),
		PublicKey:           src.PublicKey,
		EncryptedPrivateKey: src.EncryptedPrivateKey,
		Algorithm:           src.Algorithm,
		Metadata:            src.Metadata,
		Active:              src.Active,
		CreatedAt:           src.CreatedAt,
		RotatedAt:           src.RotatedAt,
		ExpiresAt:           src.ExpiresAt,
	}
}

func mapToActiveSigningKey(src sqlc.GetActiveSigningKeysRow) models.ActiveSigningKey {
	return models.ActiveSigningKey{
		ID:        src.ID,
		PublicKey: src.PublicKey,
		Algorithm: src.Algorithm,
	}
}
