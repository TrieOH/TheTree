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

var _ ports.ApiKeysRepo = (*repo)(nil)

func NewRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.ApiKeysRepo {
	return &repo{
		q:      q,
		log:    log,
		tracer: tracer,
		dbe:    database.NewErrorHandler("api keys"),
	}
}

func mapApiKey(src sqlc.ApiKey) models.ApiKey {
	return models.ApiKey{
		ID:         src.ID,
		ActorID:    src.ActorID,
		ProjectID:  src.ProjectID,
		Name:       src.Name,
		KeyPrefix:  src.KeyPrefix,
		KeyHash:    src.KeyHash,
		Metadata:   src.Metadata,
		ExpiresAt:  src.ExpiresAt,
		RevokedAt:  src.RevokedAt,
		LastUsedAt: src.LastUsedAt,
		CreatedAt:  src.CreatedAt,
	}
}
