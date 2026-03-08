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

type marketplaceConfigRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	tracer trace.Tracer
}

var _ domain.MarketplaceConfigRepo = (*marketplaceConfigRepo)(nil)

func NewMarketplaceConfigRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) domain.MarketplaceConfigRepo {
	return &marketplaceConfigRepo{q: q, log: log, tracer: tracer}
}

func (repo *marketplaceConfigRepo) queries(ctx context.Context) *sqlc.Queries {
	if tx, ok := ctx.Value(database.TxKeyValue).(pgx.Tx); ok && tx != nil {
		return repo.q.WithTx(tx)
	}
	return repo.q
}

func mapMarketplaceConfigFromDB(src *sqlc.MarketplaceConfig) *domain.MarketplaceConfig {
	return &domain.MarketplaceConfig{
		ID:           src.ID,
		WorkspaceID:  src.WorkspaceID,
		CredentialID: src.CredentialID,
		FeeBps:       src.FeeBps,
		CreatedAt:    src.CreatedAt,
		UpdatedAt:    src.UpdatedAt,
	}
}

func (repo *marketplaceConfigRepo) Create(ctx context.Context, config domain.MarketplaceConfig) (*domain.MarketplaceConfig, error) {
	ctx, span := repo.tracer.Start(ctx, "MarketplaceConfigRepo.Create")
	defer span.End()

	row, err := repo.queries(ctx).CreateMarketplaceConfig(ctx, sqlc.CreateMarketplaceConfigParams{
		WorkspaceID:  config.WorkspaceID,
		CredentialID: config.CredentialID,
		FeeBps:       config.FeeBps,
	})
	if err != nil {
		return nil, errx.FromDB(err, "marketplace_config")
	}

	return mapMarketplaceConfigFromDB(&row), nil
}

func (repo *marketplaceConfigRepo) Get(ctx context.Context, workspaceID uuid.UUID) (*domain.MarketplaceConfig, error) {
	ctx, span := repo.tracer.Start(ctx, "MarketplaceConfigRepo.Get")
	defer span.End()

	row, err := repo.queries(ctx).GetMarketplaceConfig(ctx, workspaceID)
	if err != nil {
		return nil, errx.FromDB(err, "marketplace_config")
	}

	return mapMarketplaceConfigFromDB(&row), nil
}

func (repo *marketplaceConfigRepo) Update(ctx context.Context, config domain.MarketplaceConfig) (*domain.MarketplaceConfig, error) {
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

func (repo *marketplaceConfigRepo) Delete(ctx context.Context, workspaceID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "MarketplaceConfigRepo.Delete")
	defer span.End()

	if err := repo.queries(ctx).DeleteMarketplaceConfig(ctx, workspaceID); err != nil {
		return errx.FromDB(err, "marketplace_config")
	}

	return nil
}
