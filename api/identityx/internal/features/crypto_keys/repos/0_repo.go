package repos

import (
	sqlc2 "IdentityX/internal/sqlc"
	"IdentityX/models"
	"IdentityX/ports"
	"lib/database"
)

type Repo struct {
	q   *sqlc2.Queries
	dbe database.ErrorHandler
}

var _ ports.CryptoKeysRepo = (*Repo)(nil)

func NewRepo(q *sqlc2.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler("crypto keys"),
	}
}

func mapKeys(src sqlc2.CryptoKey) models.CryptoKey {
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

func mapToActiveSigningKey(src sqlc2.GetActiveSigningKeysRow) models.ActiveSigningKey {
	return models.ActiveSigningKey{
		ID:        src.ID,
		PublicKey: src.PublicKey,
		Algorithm: src.Algorithm,
	}
}
