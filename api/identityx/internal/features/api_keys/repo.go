package api_keys

import (
	"context"

	"IdentityX/contracts"
	"IdentityX/internal/database/sqlc"
	"IdentityX/internal/shared/ports"
	"lib/database"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type apiKeyRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	tracer trace.Tracer
	dbe    database.ErrorHandler
}

var _ ports.ApiKeyRepository = (*apiKeyRepo)(nil)

func NewRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.ApiKeyRepository {
	return &apiKeyRepo{
		q:      q,
		log:    log,
		tracer: tracer,
		dbe:    database.NewErrorHandler("api key"),
	}
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
	ctx, span := repo.tracer.Start(ctx, "Upsert")
	span.SetAttributes(attribute.String("project.id", key.ProjectID.String()))
	span.SetAttributes(attribute.String("client.id", key.ClientID.String()))
	defer span.End()
	err := database.Queries(ctx, repo.q).UpsertApiKey(ctx, sqlc.UpsertApiKeyParams{
		ProjectID: key.ProjectID,
		ClientID:  key.ClientID,
		KeyHash:   key.KeyHash,
	})
	if err != nil {
		return repo.dbe(err)
	}
	return nil
}

func (repo *apiKeyRepo) GetByProjectID(ctx context.Context, projectID uuid.UUID) (*contracts.ApiKey, error) {
	ctx, span := repo.tracer.Start(ctx, "GetByProjectID")
	span.SetAttributes(attribute.String("project.id", projectID.String()))
	defer span.End()
	dbKey, err := database.Queries(ctx, repo.q).GetApiKeyByProjectID(ctx, projectID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapApiKeyFromDB(dbKey)), nil
}

func (repo *apiKeyRepo) Delete(ctx context.Context, projectID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "Delete")
	span.SetAttributes(attribute.String("project.id", projectID.String()))
	defer span.End()
	err := database.Queries(ctx, repo.q).DeleteApiKey(ctx, projectID)
	if err != nil {
		return repo.dbe(err)
	}
	return nil
}
