package collectors

import (
	"context"
	"lib/database"
	"payssage/errs"
	"payssage/internal/sqlc"
	"payssage/models"
	"payssage/ports"

	"github.com/google/uuid"
)

type Repo struct {
	q   *sqlc.Queries
	dbe database.ErrorHandler
}

var _ ports.CollectorRepo = (*Repo)(nil)

func NewRepo(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler(errs.ErrCollectorNotFound),
	}
}

// db is the queries handle, on the context's transaction when it carries one.
func (repo *Repo) db(ctx context.Context) *sqlc.Queries { return database.Queries(ctx, repo.q) }

func mapCollector(src sqlc.Collector) models.Collector {
	return models.Collector{
		ID:             src.ID,
		OwnerID:        src.OwnerID,
		OrganizationID: src.OrganizationID,
		Provider:       src.Provider,
		ProviderUserID: src.ProviderUserID,
		Credentials:    src.Credentials,
		CreatedAt:      src.CreatedAt,
		RevokedAt:      src.RevokedAt,
	}
}

func (repo *Repo) Create(ctx context.Context, toCreate models.Collector) (*models.Collector, error) {
	sqlcCollector, err := repo.db(ctx).CreateCollector(ctx, sqlc.CreateCollectorParams{
		OwnerID:        toCreate.OwnerID,
		OrganizationID: toCreate.OrganizationID,
		Provider:       toCreate.Provider,
		ProviderUserID: toCreate.ProviderUserID,
		Credentials:    toCreate.Credentials,
	})
	return database.One(sqlcCollector, err, repo.dbe, mapCollector)
}

func (repo *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Collector, error) {
	sqlcCollector, err := repo.db(ctx).GetCollectorByID(ctx, id)
	return database.One(sqlcCollector, err, repo.dbe, mapCollector)
}

func (repo *Repo) List(ctx context.Context) ([]models.Collector, error) {
	sqlcCollectors, err := repo.db(ctx).ListCollectors(ctx)
	return database.Many(sqlcCollectors, err, repo.dbe, mapCollector)
}

func (repo *Repo) ListByOrg(ctx context.Context, orgID uuid.UUID) ([]models.Collector, error) {
	sqlcCollectors, err := repo.db(ctx).ListCollectorsByOrg(ctx, &orgID)
	return database.Many(sqlcCollectors, err, repo.dbe, mapCollector)
}

func (repo *Repo) ListByOwner(ctx context.Context, ownerID uuid.UUID) ([]models.Collector, error) {
	sqlcCollectors, err := repo.db(ctx).ListCollectorsByOwner(ctx, ownerID)
	return database.Many(sqlcCollectors, err, repo.dbe, mapCollector)
}

func (repo *Repo) Revoke(ctx context.Context, id uuid.UUID) error {
	return repo.dbe(repo.db(ctx).RevokeCollector(ctx, id))
}
