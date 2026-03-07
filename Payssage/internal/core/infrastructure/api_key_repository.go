package infrastructure

import (
	"TriePayments/internal/core/domain"
	"TriePayments/internal/plataform/database"
	"TriePayments/internal/plataform/database/sqlc"
	"TriePayments/internal/shared/errx"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type apiKeyRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	tracer trace.Tracer
}

var _ domain.ApiKeysRepo = (*apiKeyRepo)(nil)

func NewApiKeyRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) domain.ApiKeysRepo {
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

func mapApiKeyFromDB(src *sqlc.ApiKey) *domain.APIKey {
	return &domain.APIKey{
		ID:          src.ID,
		ScopeID:     src.ScopeID,
		WorkspaceID: src.WorkspaceID,
		Name:        src.Name,
		KeyHash:     src.KeyHash,
		KeyPrefix:   src.KeyPrefix,
		CreatedAt:   src.CreatedAt,
		RevokedAt:   src.RevokedAt,
	}
}

func (repo *apiKeyRepo) Create(ctx context.Context, toCreate domain.APIKey) (*domain.APIKey, error) {
	ctx, span := repo.tracer.Start(ctx, "ApiKeyRepo.Create")
	defer span.End()

	sqlcApiKey, err := repo.queries(ctx).CreateAPIKey(ctx, sqlc.CreateAPIKeyParams{
		ID:          toCreate.ID,
		ScopeID:     toCreate.ScopeID,
		WorkspaceID: toCreate.WorkspaceID,
		Name:        toCreate.Name,
		KeyHash:     toCreate.KeyHash,
		KeyPrefix:   toCreate.KeyPrefix,
	})
	if err != nil {
		return nil, errx.FromDB(err, "api key")
	}

	return mapApiKeyFromDB(&sqlcApiKey), nil
}

func (repo *apiKeyRepo) GetByPrefix(ctx context.Context, prefix string) ([]domain.APIKey, error) {
	ctx, span := repo.tracer.Start(ctx, "WorkspaceRepo.GetByPrefix")
	defer span.End()

	sqlcApiKeys, err := repo.queries(ctx).GetAPIKeyByPrefix(ctx, prefix)
	if err != nil {
		return nil, errx.FromDB(err, "api key")
	}

	out := make([]domain.APIKey, 0, len(sqlcApiKeys))
	for _, key := range sqlcApiKeys {
		out = append(out, *mapApiKeyFromDB(&key))
	}
	return out, nil
}

func (repo *apiKeyRepo) ListByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]domain.APIKey, error) {
	ctx, span := repo.tracer.Start(ctx, "WorkspaceRepo.ListByWorkspace")
	defer span.End()

	sqlcApiKeys, err := repo.queries(ctx).ListAPIKeysByWorkspace(ctx, workspaceID)
	if err != nil {
		return nil, errx.FromDB(err, "api key")
	}

	out := make([]domain.APIKey, 0, len(sqlcApiKeys))
	for _, key := range sqlcApiKeys {
		out = append(out, *mapApiKeyFromDB(&key))
	}
	return out, nil
}

func (repo *apiKeyRepo) Revoke(ctx context.Context, id, workspaceID uuid.UUID) (*domain.APIKey, error) {
	ctx, span := repo.tracer.Start(ctx, "ApiKeyRepo.Revoke")
	defer span.End()

	sqlcApiKey, err := repo.queries(ctx).RevokeAPIKey(ctx, sqlc.RevokeAPIKeyParams{
		ID:          id,
		WorkspaceID: workspaceID,
	})
	if err != nil {
		return nil, errx.FromDB(err, "api key")
	}

	return mapApiKeyFromDB(&sqlcApiKey), nil
}
