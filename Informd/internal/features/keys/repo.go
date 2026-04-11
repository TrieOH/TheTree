package keys

import (
	"TrieForms/internal/platform/database"
	"TrieForms/internal/platform/database/sqlc"
	"TrieForms/internal/shared/contracts"
	"TrieForms/internal/shared/errx"
	"TrieForms/internal/shared/ports"
	"context"

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

func NewApiKeyRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.ApiKeysRepo {
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

func mapApiKeyFromDB(src *sqlc.ApiKey) *contracts.APIKey {
	return &contracts.APIKey{
		ID:        src.ID,
		OwnerID:   src.OwnerID,
		ScopeID:   src.ScopeID,
		ProjectID: src.ProjectID,
		Name:      src.Name,
		KeyHash:   src.KeyHash,
		KeyPrefix: src.KeyPrefix,
		CreatedAt: src.CreatedAt,
		RevokedAt: src.RevokedAt,
	}
}

func (repo *repo) Create(ctx context.Context, toCreate contracts.APIKey) (*contracts.APIKey, error) {
	ctx, span := repo.tracer.Start(ctx, "ApiKeyRepo.Create")
	defer span.End()

	sqlcApiKey, err := repo.queries(ctx).CreateAPIKey(ctx, sqlc.CreateAPIKeyParams{
		ID:        toCreate.ID,
		OwnerID:   toCreate.OwnerID,
		ScopeID:   toCreate.ScopeID,
		ProjectID: toCreate.ProjectID,
		Name:      toCreate.Name,
		KeyHash:   toCreate.KeyHash,
		KeyPrefix: toCreate.KeyPrefix,
	})
	if err != nil {
		return nil, errx.FromDB(err, "api key")
	}

	return mapApiKeyFromDB(&sqlcApiKey), nil
}

func (repo *repo) GetByPrefix(ctx context.Context, prefix string) ([]contracts.APIKey, error) {
	ctx, span := repo.tracer.Start(ctx, "ApiKeyRepo.Create")
	defer span.End()

	sqlcApiKeys, err := repo.queries(ctx).GetAPIKeyByPrefix(ctx, prefix)
	if err != nil {
		return nil, errx.FromDB(err, "api key")
	}

	out := make([]contracts.APIKey, 0, len(sqlcApiKeys))
	for _, key := range sqlcApiKeys {
		out = append(out, *mapApiKeyFromDB(&key))
	}
	return out, nil
}

func (repo *repo) ListByProject(ctx context.Context, projectID uuid.UUID) ([]contracts.APIKey, error) {
	ctx, span := repo.tracer.Start(ctx, "ApiKeyRepo.ListByProject")
	defer span.End()

	sqlcApiKeys, err := repo.queries(ctx).ListAPIKeysByProject(ctx, projectID)
	if err != nil {
		return nil, errx.FromDB(err, "api key")
	}

	out := make([]contracts.APIKey, 0, len(sqlcApiKeys))
	for _, key := range sqlcApiKeys {
		out = append(out, *mapApiKeyFromDB(&key))
	}
	return out, nil
}

func (repo *repo) Revoke(ctx context.Context, id, userID uuid.UUID) (*contracts.APIKey, error) {
	ctx, span := repo.tracer.Start(ctx, "ApiKeyRepo.Revoke")
	defer span.End()

	sqlcApiKey, err := repo.queries(ctx).RevokeAPIKey(ctx, sqlc.RevokeAPIKeyParams{
		ID:      id,
		OwnerID: userID,
	})
	if err != nil {
		return nil, errx.FromDB(err, "api key")
	}

	return mapApiKeyFromDB(&sqlcApiKey), nil
}
