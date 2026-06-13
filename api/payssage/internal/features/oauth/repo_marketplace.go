package oauth

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

type marketplaceConfigRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	tracer trace.Tracer
}

var _ ports.MarketplaceConfigRepo = (*marketplaceConfigRepo)(nil)

func NewMarketplaceConfigRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.MarketplaceConfigRepo {
	return &marketplaceConfigRepo{q: q, log: log, tracer: tracer}
}

func (repo *marketplaceConfigRepo) queries(ctx context.Context) *sqlc.Queries {
	if tx, ok := ctx.Value(database.TxKeyValue).(pgx.Tx); ok && tx != nil {
		return repo.q.WithTx(tx)
	}
	return repo.q
}

func mapMarketplaceConfigFromDB(src *sqlc.MarketplaceConfig) *models.MarketplaceConfig {
	return &models.MarketplaceConfig{
		ID:           src.ID,
		WorkspaceID:  src.WorkspaceID,
		CredentialID: src.CredentialID,
		Provider:     src.Provider,
		FeeBps:       src.FeeBps,
		CreatedAt:    src.CreatedAt,
		UpdatedAt:    src.UpdatedAt,
	}
}

func (repo *marketplaceConfigRepo) Create(ctx context.Context, config models.MarketplaceConfig) (*models.MarketplaceConfig, error) {
	ctx, span := repo.tracer.Start(ctx, "MarketplaceConfigRepo.Create")
	defer span.End()

	row, err := repo.queries(ctx).CreateMarketplaceConfig(ctx, sqlc.CreateMarketplaceConfigParams{
		WorkspaceID:  config.WorkspaceID,
		CredentialID: config.CredentialID,
		Provider:     config.Provider,
		FeeBps:       config.FeeBps,
	})
	if err != nil {
		return nil, errx.FromDB(err, "marketplace_config")
	}

	return mapMarketplaceConfigFromDB(&row), nil
}

func (repo *marketplaceConfigRepo) List(ctx context.Context, workspaceID uuid.UUID) ([]models.MarketplaceConfig, error) {
	ctx, span := repo.tracer.Start(ctx, "MarketplaceConfigRepo.List")
	defer span.End()

	rows, err := repo.queries(ctx).ListMarketplaceConfigs(ctx, workspaceID)
	if err != nil {
		return nil, errx.FromDB(err, "marketplace_config")
	}

	configs := make([]models.MarketplaceConfig, len(rows))
	for i := range rows {
		configs[i] = *mapMarketplaceConfigFromDB(&rows[i])
	}
	return configs, nil
}

func (repo *marketplaceConfigRepo) Get(ctx context.Context, workspaceID, credentialID uuid.UUID) (*models.MarketplaceConfig, error) {
	ctx, span := repo.tracer.Start(ctx, "MarketplaceConfigRepo.Get")
	defer span.End()

	row, err := repo.queries(ctx).GetMarketplaceConfig(ctx, sqlc.GetMarketplaceConfigParams{
		WorkspaceID:  workspaceID,
		CredentialID: credentialID,
	})
	if err != nil {
		return nil, errx.FromDB(err, "marketplace_config")
	}

	return mapMarketplaceConfigFromDB(&row), nil
}

func (repo *marketplaceConfigRepo) GetByProvider(ctx context.Context, workspaceID uuid.UUID, provider string) (*models.MarketplaceConfig, error) {
	ctx, span := repo.tracer.Start(ctx, "MarketplaceConfigRepo.GetByProvider")
	defer span.End()

	row, err := repo.queries(ctx).GetMarketplaceConfigByProvider(ctx, sqlc.GetMarketplaceConfigByProviderParams{
		WorkspaceID: workspaceID,
		Provider:    provider,
	})
	if err != nil {
		return nil, errx.FromDB(err, "marketplace_config")
	}

	return mapMarketplaceConfigFromDB(&row), nil
}

func (repo *marketplaceConfigRepo) Update(ctx context.Context, config models.MarketplaceConfig) (*models.MarketplaceConfig, error) {
	ctx, span := repo.tracer.Start(ctx, "MarketplaceConfigRepo.Update")
	defer span.End()

	row, err := repo.queries(ctx).UpdateMarketplaceConfig(ctx, sqlc.UpdateMarketplaceConfigParams{
		WorkspaceID:  config.WorkspaceID,
		CredentialID: config.CredentialID,
		FeeBps:       config.FeeBps,
	})
	if err != nil {
		return nil, errx.FromDB(err, "marketplace_config")
	}

	return mapMarketplaceConfigFromDB(&row), nil
}

func (repo *marketplaceConfigRepo) Delete(ctx context.Context, workspaceID, credentialID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "MarketplaceConfigRepo.Delete")
	defer span.End()

	if err := repo.queries(ctx).DeleteMarketplaceConfig(ctx, sqlc.DeleteMarketplaceConfigParams{
		WorkspaceID:  workspaceID,
		CredentialID: credentialID,
	}); err != nil {
		return errx.FromDB(err, "marketplace_config")
	}

	return nil
}

func (repo *marketplaceConfigRepo) DeleteAll(ctx context.Context, workspaceID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "MarketplaceConfigRepo.DeleteAll")
	defer span.End()

	if err := repo.queries(ctx).DeleteAllMarketplaceConfigs(ctx, workspaceID); err != nil {
		return errx.FromDB(err, "marketplace_config")
	}

	return nil
}
