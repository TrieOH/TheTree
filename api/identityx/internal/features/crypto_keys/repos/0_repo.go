package repos

import (
	"IdentityX/internal/database/sqlc"
	"IdentityX/models"
	"IdentityX/ports"
	"lib/database"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type repo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	tracer trace.Tracer
	dbe    database.ErrorHandler
}

var _ ports.CryptoKeysRepo = (*repo)(nil)

func NewRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.CryptoKeysRepo {
	return &repo{
		q:      q,
		log:    log,
		tracer: tracer,
		dbe:    database.NewErrorHandler("crypto keys"),
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
