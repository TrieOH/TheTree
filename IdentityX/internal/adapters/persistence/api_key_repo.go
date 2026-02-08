package persistence

import (
	"GoAuth/internal/adapters/persistence/sqlc"
	"GoAuth/internal/adapters/persistence/transactions"
	"GoAuth/internal/domain/apikey"
	"GoAuth/internal/ports/outbounds"
	"context"

	"github.com/MintzyG/fail/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type apiKeyRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	tracer trace.Tracer
}

var _ outbounds.ApiKeyRepository = (*apiKeyRepo)(nil)

func NewApiKeyRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) outbounds.ApiKeyRepository {
	return &apiKeyRepo{
		q:      q,
		log:    log,
		tracer: tracer,
	}
}

func (repo *apiKeyRepo) queries(ctx context.Context) *sqlc.Queries {
	if tx, ok := ctx.Value(transactions.TxKeyValue).(pgx.Tx); ok && tx != nil {
		return repo.q.WithTx(tx)
	}
	return repo.q
}

func mapApiKeyFromDB(dst *apikey.ApiKey, src *sqlc.ApiKey) {
	dst.ProjectID = src.ProjectID
	dst.ClientID = src.ClientID
	dst.KeyHash = src.KeyHash
	dst.CreatedAt = src.CreatedAt
	dst.UpdatedAt = src.UpdatedAt
}

func (repo *apiKeyRepo) Upsert(ctx context.Context, key apikey.ApiKey) error {
	ctx, span := repo.tracer.Start(ctx, "ApiKeyRepo.Upsert",
		trace.WithAttributes(
			attribute.String("project.id", key.ProjectID.String()),
			attribute.String("client.id", key.ClientID.String()),
		),
	)
	defer span.End()

	err := repo.queries(ctx).UpsertApiKey(ctx, sqlc.UpsertApiKeyParams{
		ProjectID: key.ProjectID,
		ClientID:  key.ClientID,
		KeyHash:   key.KeyHash,
	})
	if err != nil {
		return fail.From(err).RecordCtx(ctx)
	}

	return nil
}

func (repo *apiKeyRepo) GetByProjectID(ctx context.Context, projectID uuid.UUID) (*apikey.ApiKey, error) {
	ctx, span := repo.tracer.Start(ctx, "ApiKeyRepo.GetByProjectID",
		trace.WithAttributes(
			attribute.String("project.id", projectID.String()),
		),
	)
	defer span.End()

	dbKey, err := repo.queries(ctx).GetApiKeyByProjectID(ctx, projectID)
	if err != nil {
		return nil, fail.From(err).WithArgs("api_key").RecordCtx(ctx)
	}

	var key apikey.ApiKey
	mapApiKeyFromDB(&key, &dbKey)
	return &key, nil
}

func (repo *apiKeyRepo) Delete(ctx context.Context, projectID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "ApiKeyRepo.Delete",
		trace.WithAttributes(
			attribute.String("project.id", projectID.String()),
		),
	)
	defer span.End()

	err := repo.queries(ctx).DeleteApiKey(ctx, projectID)
	if err != nil {
		return fail.From(err).RecordCtx(ctx)
	}

	return nil
}
