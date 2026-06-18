package oauth

import (
	"context"
	"encoding/json"
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

type providerCredentialsRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	tracer trace.Tracer
}

var _ ports.ProviderCredentialRepo = (*providerCredentialsRepo)(nil)

func NewProviderCredentialsRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.ProviderCredentialRepo {
	return &providerCredentialsRepo{q: q, log: log, tracer: tracer}
}

func (repo *providerCredentialsRepo) queries(ctx context.Context) *sqlc.Queries {
	if tx, ok := ctx.Value(database.TxKeyValue).(pgx.Tx); ok && tx != nil {
		return repo.q.WithTx(tx)
	}
	return repo.q
}

func mapProviderCredentialFromDB(src *sqlc.ProviderCredential) (*models.ProviderCredential, error) {
	var data models.ProviderCredentialData
	if err := json.Unmarshal(src.Credentials, &data); err != nil {
		return nil, err
	}

	return &models.ProviderCredential{
		ID:          src.ID,
		WorkspaceID: src.WorkspaceID,
		Provider:    src.Provider,
		Credentials: data,
		CreatedAt:   src.CreatedAt,
		RevokedAt:   src.RevokedAt,
	}, nil
}

func (repo *providerCredentialsRepo) Create(ctx context.Context, cred models.ProviderCredential) (*models.ProviderCredential, error) {
	ctx, span := repo.tracer.Start(ctx, "ProviderCredentialsRepo.Create")
	defer span.End()

	credJSON, err := json.Marshal(cred.Credentials)
	if err != nil {
		return nil, errx.Internal("provider_credential").SetCause(err)
	}

	row, err := repo.queries(ctx).CreateProviderCredential(ctx, sqlc.CreateProviderCredentialParams{
		WorkspaceID: cred.WorkspaceID,
		Provider:    cred.Provider,
		Credentials: credJSON,
	})
	if err != nil {
		return nil, errx.FromDB(err, "provider_credential")
	}

	return mapProviderCredentialFromDB(&row)
}

func (repo *providerCredentialsRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.ProviderCredential, error) {
	ctx, span := repo.tracer.Start(ctx, "ProviderCredentialsRepo.GetByID")
	defer span.End()

	row, err := repo.queries(ctx).GetProviderCredential(ctx, id)
	if err != nil {
		return nil, errx.FromDB(err, "provider_credential")
	}

	return mapProviderCredentialFromDB(&row)
}

func (repo *providerCredentialsRepo) ListByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]models.ProviderCredential, error) {
	ctx, span := repo.tracer.Start(ctx, "ProviderCredentialsRepo.ListByWorkspace")
	defer span.End()

	rows, err := repo.queries(ctx).ListProviderCredentials(ctx, workspaceID)
	if err != nil {
		return nil, errx.FromDB(err, "provider_credential")
	}

	out := make([]models.ProviderCredential, 0, len(rows))
	for _, row := range rows {
		cred, err := mapProviderCredentialFromDB(&row)
		if err != nil {
			return nil, errx.Internal("provider_credential").SetCause(err)
		}
		out = append(out, *cred)
	}
	return out, nil
}

func (repo *providerCredentialsRepo) Revoke(ctx context.Context, id uuid.UUID, workspaceID uuid.UUID) (*models.ProviderCredential, error) {
	ctx, span := repo.tracer.Start(ctx, "ProviderCredentialsRepo.Revoke")
	defer span.End()

	row, err := repo.queries(ctx).RevokeProviderCredential(ctx, sqlc.RevokeProviderCredentialParams{
		ID:          id,
		WorkspaceID: workspaceID,
	})
	if err != nil {
		return nil, errx.FromDB(err, "provider_credential")
	}

	return mapProviderCredentialFromDB(&row)
}

func (repo *providerCredentialsRepo) GetByWorkspaceAndProvider(ctx context.Context, workspaceID uuid.UUID, provider string) (*models.ProviderCredential, error) {
	ctx, span := repo.tracer.Start(ctx, "ProviderCredentialsRepo.GetByWorkspaceAndProvider")
	defer span.End()

	row, err := repo.queries(ctx).GetWorkspaceProviderCredential(ctx, sqlc.GetWorkspaceProviderCredentialParams{
		WorkspaceID: workspaceID,
		Provider:    provider,
	})
	if err != nil {
		return nil, errx.FromDB(err, "provider_credential")
	}

	return mapProviderCredentialFromDB(&row)
}

func (repo *providerCredentialsRepo) GetSellerCredentialByProvider(ctx context.Context, workspaceID uuid.UUID, provider string) (*models.ProviderCredential, error) {
	ctx, span := repo.tracer.Start(ctx, "ProviderCredentialsRepo.GetSellerCredentialByProvider")
	defer span.End()

	row, err := repo.queries(ctx).GetSellerCredentialByProvider(ctx, sqlc.GetSellerCredentialByProviderParams{
		WorkspaceID: workspaceID,
		Provider:    provider,
	})
	if err != nil {
		return nil, errx.FromDB(err, "provider_credential")
	}
	return mapProviderCredentialFromDB(&row)
}
