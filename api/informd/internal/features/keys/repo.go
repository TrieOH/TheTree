package keys

import (
	"Informd/internal/database/sqlc"
	"Informd/models"
	"Informd/ports"
	"context"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
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

func (repo *repo) Create(ctx context.Context, toCreate models.APIKey) (*models.APIKey, error) {
	ctx, span := repo.tracer.Start(ctx, "Create")
	defer span.End()
	sqlcApiKey, err := database.Queries(ctx, repo.q).CreateAPIKey(ctx, sqlc.CreateAPIKeyParams{
		OwnerID:   toCreate.OwnerID,
		Name:      toCreate.Name,
		KeyHash:   toCreate.KeyHash,
		KeyPrefix: toCreate.KeyPrefix,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapApiKey(sqlcApiKey)), nil
}

func (repo *repo) GetByPrefix(ctx context.Context, prefix string) ([]models.APIKey, error) {
	ctx, span := repo.tracer.Start(ctx, "GetByPrefix")
	defer span.End()
	sqlcApiKeys, err := database.Queries(ctx, repo.q).GetAPIKeyByPrefix(ctx, prefix)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcApiKeys, mapApiKey), nil
}

func (repo *repo) Revoke(ctx context.Context, id, userID uuid.UUID) (*models.APIKey, error) {
	ctx, span := repo.tracer.Start(ctx, "Revoke")
	defer span.End()
	sqlcApiKey, err := database.Queries(ctx, repo.q).RevokeAPIKey(ctx, sqlc.RevokeAPIKeyParams{
		ID:      id,
		OwnerID: userID,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapApiKey(sqlcApiKey)), nil
}

func (repo *repo) BulkGet(ctx context.Context, ids []uuid.UUID) ([]models.APIKey, error) {
	ctx, span := repo.tracer.Start(ctx, "BulkGet")
	defer span.End()
	sqlcKeys, err := database.Queries(ctx, repo.q).BulkGetAPIKeys(ctx, ids)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcKeys, mapApiKey), nil
}
