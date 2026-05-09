package keys

import (
	"Informd/internal/platform/database/sqlc"
	"Informd/internal/shared/contracts"
	"Informd/internal/shared/errx"
	"Informd/internal/shared/ports"
	"context"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type repo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	tracer trace.Tracer
}

var _ ports.ApiKeysRepo = (*repo)(nil)

func NewRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.ApiKeysRepo {
	return &repo{
		q:      q,
		log:    log,
		tracer: tracer,
	}
}

func (repo *repo) queries(ctx context.Context) *sqlc.Queries {
	if tx, ok := ctx.Value(database.TxKeyValue).(pgx.Tx); ok && tx != nil {
		return repo.q.WithTx(tx)
	}
	return repo.q
}

func (repo *repo) span(ctx context.Context, op string) (context.Context, trace.Span) {
	return repo.tracer.Start(ctx, "ApiKeyRepo."+op)
}

func mapApiKey(src sqlc.ApiKey) contracts.APIKey {
	return contracts.APIKey{
		ID:        src.ID,
		OwnerID:   src.OwnerID,
		Name:      src.Name,
		KeyHash:   src.KeyHash,
		KeyPrefix: src.KeyPrefix,
		CreatedAt: src.CreatedAt,
		RevokedAt: src.RevokedAt,
	}
}

func (repo *repo) Create(ctx context.Context, toCreate contracts.APIKey) (*contracts.APIKey, error) {
	ctx, span := repo.span(ctx, "Create")
	defer span.End()
	sqlcApiKey, err := repo.queries(ctx).CreateAPIKey(ctx, sqlc.CreateAPIKeyParams{
		OwnerID:   toCreate.OwnerID,
		Name:      toCreate.Name,
		KeyHash:   toCreate.KeyHash,
		KeyPrefix: toCreate.KeyPrefix,
	})
	if err != nil {
		return nil, errx.DB(err, "api key")
	}
	return new(mapApiKey(sqlcApiKey)), nil
}

func (repo *repo) GetByPrefix(ctx context.Context, prefix string) ([]contracts.APIKey, error) {
	ctx, span := repo.span(ctx, "GetByPrefix")
	defer span.End()
	sqlcApiKeys, err := repo.queries(ctx).GetAPIKeyByPrefix(ctx, prefix)
	if err != nil {
		return nil, errx.DB(err, "api key")
	}
	return xslices.MapSlice(sqlcApiKeys, mapApiKey), nil
}

func (repo *repo) Revoke(ctx context.Context, id, userID uuid.UUID) (*contracts.APIKey, error) {
	ctx, span := repo.span(ctx, "Revoke")
	defer span.End()
	sqlcApiKey, err := repo.queries(ctx).RevokeAPIKey(ctx, sqlc.RevokeAPIKeyParams{
		ID:      id,
		OwnerID: userID,
	})
	if err != nil {
		return nil, errx.DB(err, "api key")
	}
	return new(mapApiKey(sqlcApiKey)), nil
}

func (repo *repo) BulkGet(ctx context.Context, ids []uuid.UUID) ([]contracts.APIKey, error) {
	ctx, span := repo.span(ctx, "BulkGet")
	defer span.End()
	sqlcKeys, err := repo.queries(ctx).BulkGetAPIKeys(ctx, ids)
	if err != nil {
		return nil, errx.DB(err, "api key")
	}
	return xslices.MapSlice(sqlcKeys, mapApiKey), nil
}
