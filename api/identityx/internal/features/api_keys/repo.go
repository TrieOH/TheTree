package api_keys

import (
	"IdentityX/contracts"
	"IdentityX/internal/platform/database/sqlc"
	"IdentityX/internal/shared/ports"
	"context"
	"lib/database"
	"lib/errx"

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
	dbe    *errx.DBHandler
}

var _ ports.ApiKeyRepository = (*apiKeyRepo)(nil)

func NewRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer, dbe *errx.DBHandler) ports.ApiKeyRepository {
	return &apiKeyRepo{
		q:      q,
		log:    log,
		tracer: tracer,
		dbe:    dbe,
	}
}

func (repo *apiKeyRepo) queries(ctx context.Context) *sqlc.Queries {
	if tx, ok := ctx.Value(database.TxKeyValue).(pgx.Tx); ok && tx != nil {
		return repo.q.WithTx(tx)
	}
	return repo.q
}

func (repo *apiKeyRepo) span(ctx context.Context, op string) (context.Context, trace.Span) {
	return repo.tracer.Start(ctx, "ApiKeyRepo."+op)
}

func mapApiKeyFromDB(src sqlc.ApiKey) contracts.ApiKey {
	return contracts.ApiKey{
		ProjectID: src.ProjectID,
		ClientID:  src.ClientID,
		KeyHash:   src.KeyHash,
		CreatedAt: src.CreatedAt,
		UpdatedAt: src.UpdatedAt,
	}

}

func (repo *apiKeyRepo) Upsert(ctx context.Context, key contracts.ApiKey) error {
	ctx, span := repo.span(ctx, "Upsert")
	span.SetAttributes(attribute.String("project.id", key.ProjectID.String()))
	span.SetAttributes(attribute.String("client.id", key.ClientID.String()))
	defer span.End()
	err := repo.queries(ctx).UpsertApiKey(ctx, sqlc.UpsertApiKeyParams{
		ProjectID: key.ProjectID,
		ClientID:  key.ClientID,
		KeyHash:   key.KeyHash,
	})
	if err != nil {
		return repo.dbe.DB(err, "api key")
	}
	return nil
}

func (repo *apiKeyRepo) GetByProjectID(ctx context.Context, projectID uuid.UUID) (*contracts.ApiKey, error) {
	ctx, span := repo.span(ctx, "GetByProjectID")
	span.SetAttributes(attribute.String("project.id", projectID.String()))
	defer span.End()
	dbKey, err := repo.queries(ctx).GetApiKeyByProjectID(ctx, projectID)
	if err != nil {
		return nil, repo.dbe.DB(err, "api key")
	}
	return new(mapApiKeyFromDB(dbKey)), nil
}

func (repo *apiKeyRepo) Delete(ctx context.Context, projectID uuid.UUID) error {
	ctx, span := repo.span(ctx, "Delete")
	span.SetAttributes(attribute.String("project.id", projectID.String()))
	defer span.End()
	err := repo.queries(ctx).DeleteApiKey(ctx, projectID)
	if err != nil {
		return repo.dbe.DB(err, "api key")
	}
	return nil
}
