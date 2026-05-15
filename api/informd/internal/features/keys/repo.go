package keys

import (
	sqlc2 "Informd/internal/database/sqlc"
	"Informd/models"
	"Informd/ports"
	"context"
	"lib/database"
	"lib/errx"
	"lib/xslices"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type repo struct {
	q      *sqlc2.Queries
	log    *zap.Logger
	tracer trace.Tracer
	dbe    *errx.DBHandler
}

var _ ports.ApiKeysRepo = (*repo)(nil)

func NewRepo(q *sqlc2.Queries, log *zap.Logger, tracer trace.Tracer, dbe *errx.DBHandler) ports.ApiKeysRepo {
	return &repo{
		q:      q,
		log:    log,
		tracer: tracer,
		dbe:    dbe,
	}
}

func (repo *repo) queries(ctx context.Context) *sqlc2.Queries {
	if tx, ok := ctx.Value(database.TxKeyValue).(pgx.Tx); ok && tx != nil {
		return repo.q.WithTx(tx)
	}
	return repo.q
}

func (repo *repo) span(ctx context.Context, op string) (context.Context, trace.Span) {
	return repo.tracer.Start(ctx, "ApiKeyRepo."+op)
}

func mapApiKey(src sqlc2.ApiKey) models.APIKey {
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
	ctx, span := repo.span(ctx, "Create")
	defer span.End()
	sqlcApiKey, err := repo.queries(ctx).CreateAPIKey(ctx, sqlc2.CreateAPIKeyParams{
		OwnerID:   toCreate.OwnerID,
		Name:      toCreate.Name,
		KeyHash:   toCreate.KeyHash,
		KeyPrefix: toCreate.KeyPrefix,
	})
	if err != nil {
		return nil, repo.dbe.DB(err, "api key")
	}
	return new(mapApiKey(sqlcApiKey)), nil
}

func (repo *repo) GetByPrefix(ctx context.Context, prefix string) ([]models.APIKey, error) {
	ctx, span := repo.span(ctx, "GetByPrefix")
	defer span.End()
	sqlcApiKeys, err := repo.queries(ctx).GetAPIKeyByPrefix(ctx, prefix)
	if err != nil {
		return nil, repo.dbe.DB(err, "api key")
	}
	return xslices.MapSlice(sqlcApiKeys, mapApiKey), nil
}

func (repo *repo) Revoke(ctx context.Context, id, userID uuid.UUID) (*models.APIKey, error) {
	ctx, span := repo.span(ctx, "Revoke")
	defer span.End()
	sqlcApiKey, err := repo.queries(ctx).RevokeAPIKey(ctx, sqlc2.RevokeAPIKeyParams{
		ID:      id,
		OwnerID: userID,
	})
	if err != nil {
		return nil, repo.dbe.DB(err, "api key")
	}
	return new(mapApiKey(sqlcApiKey)), nil
}

func (repo *repo) BulkGet(ctx context.Context, ids []uuid.UUID) ([]models.APIKey, error) {
	ctx, span := repo.span(ctx, "BulkGet")
	defer span.End()
	sqlcKeys, err := repo.queries(ctx).BulkGetAPIKeys(ctx, ids)
	if err != nil {
		return nil, repo.dbe.DB(err, "api key")
	}
	return xslices.MapSlice(sqlcKeys, mapApiKey), nil
}
