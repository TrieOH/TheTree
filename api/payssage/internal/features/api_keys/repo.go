package api_keys

import (
	"context"
	"payssage/ports"

	"lib/database"
	"payssage/internal/database/sqlc"
	"payssage/internal/shared/errx"
	"payssage/models"

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

var _ ports.ApiKeysRepo = (*apiKeyRepo)(nil)

func NewApiKeyRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.ApiKeysRepo {
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

func mapApiKeyFromDB(src *sqlc.ApiKey) *models.APIKey {
	return &models.APIKey{
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

func (repo *apiKeyRepo) Create(ctx context.Context, toCreate models.APIKey) (*models.APIKey, error) {
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

func (repo *apiKeyRepo) GetByPrefix(ctx context.Context, prefix string) ([]models.APIKey, error) {
	ctx, span := repo.tracer.Start(ctx, "WorkspaceRepo.GetByPrefix")
	defer span.End()

	sqlcApiKeys, err := repo.queries(ctx).GetAPIKeyByPrefix(ctx, prefix)
	if err != nil {
		return nil, errx.FromDB(err, "api key")
	}

	out := make([]models.APIKey, 0, len(sqlcApiKeys))
	for _, key := range sqlcApiKeys {
		out = append(out, *mapApiKeyFromDB(&key))
	}
	return out, nil
}

func (repo *apiKeyRepo) ListByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]models.APIKey, error) {
	ctx, span := repo.tracer.Start(ctx, "WorkspaceRepo.ListByWorkspace")
	defer span.End()

	sqlcApiKeys, err := repo.queries(ctx).ListAPIKeysByWorkspace(ctx, workspaceID)
	if err != nil {
		return nil, errx.FromDB(err, "api key")
	}

	out := make([]models.APIKey, 0, len(sqlcApiKeys))
	for _, key := range sqlcApiKeys {
		out = append(out, *mapApiKeyFromDB(&key))
	}
	return out, nil
}

func (repo *apiKeyRepo) Revoke(ctx context.Context, id, workspaceID uuid.UUID) (*models.APIKey, error) {
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
