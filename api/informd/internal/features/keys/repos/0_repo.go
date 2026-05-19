package repos

import (
	"Informd/internal/database/sqlc"
	"Informd/models"
	"Informd/ports"
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

func NewRepos(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.ApiKeysRepo {
	return &repo{
		q:      q,
		log:    log,
		tracer: tracer,
		dbe:    database.NewErrorHandler("api key"),
	}
}

func mapApiKey(src sqlc.ApiKey) models.APIKey {
	return models.APIKey{
		ID:        src.ID,
		OwnerID:   src.OwnerID,
		Name:      src.Name,
		KeyHash:   src.KeyHash,
		KeyPrefix: src.KeyPrefix,
		CreatedAt: src.CreatedAt,
		RevokedAt: src.RevokedAt,
	}
}
