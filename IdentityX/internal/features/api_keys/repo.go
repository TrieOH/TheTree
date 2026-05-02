package api_keys

import (
	"IdentityX/internal/platform/database"
	"IdentityX/internal/platform/database/sqlc"
	"IdentityX/internal/shared/contracts"
	"IdentityX/internal/shared/errx"
	"IdentityX/internal/shared/ports"
	"context"

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

var _ ports.ApiKeyRepository = (*apiKeyRepo)(nil)

func NewRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.ApiKeyRepository {
	return &apiKeyRepo{
		q:      q,
		log:    log,
		tracer: tracer,
	}
}

func (repo *apiKeyRepo) queries(ctx context.Context) *sqlc.Queries {
	if tx, ok := ctx.Value(database.TxKeyValue).(pgx.Tx); ok && tx != nil {
		return repo.q.WithTx(tx)
	}
	return repo.q
}

func mapApiKeyFromDB(dst *contracts.ApiKey, src *sqlc.ApiKey) {
	dst.ProjectID = src.ProjectID
	dst.ClientID = src.ClientID
	dst.KeyHash = src.KeyHash
	dst.CreatedAt = src.CreatedAt
	dst.UpdatedAt = src.UpdatedAt
}

func (repo *apiKeyRepo) Upsert(ctx context.Context, key contracts.ApiKey) error {
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
		return errx.DB(err, "api key")
	}

	return nil
}

func (repo *apiKeyRepo) GetByProjectID(ctx context.Context, projectID uuid.UUID) (*contracts.ApiKey, error) {
	ctx, span := repo.tracer.Start(ctx, "ApiKeyRepo.GetByProjectID",
		trace.WithAttributes(
			attribute.String("project.id", projectID.String()),
		),
	)
	defer span.End()

	dbKey, err := repo.queries(ctx).GetApiKeyByProjectID(ctx, projectID)
	if err != nil {
		return nil, errx.DB(err, "api key")
	}

	var key contracts.ApiKey
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
		return errx.DB(err, "api key")
	}

	return nil
}
