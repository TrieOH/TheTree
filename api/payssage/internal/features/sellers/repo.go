package sellers

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

var _ ports.SellerRepo = (*Repo)(nil)

func NewRepo(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler(errs.ErrSellerNotFound),
	}
}

// db is the queries handle, on the context's transaction when it carries one.
func (repo *Repo) db(ctx context.Context) *sqlc.Queries { return database.Queries(ctx, repo.q) }

func mapSeller(src sqlc.Seller) models.Seller {
	return models.Seller{
		ID:             src.ID,
		WalletID:       src.WalletID,
		Provider:       src.Provider,
		ProviderUserID: src.ProviderUserID,
		Credentials:    src.Credentials,
		CreatedAt:      src.CreatedAt,
		RevokedAt:      src.RevokedAt,
	}
}

func (repo *Repo) Create(ctx context.Context, toCreate models.Seller) (*models.Seller, error) {
	sqlcSeller, err := repo.db(ctx).CreateSeller(ctx, sqlc.CreateSellerParams{
		WalletID:       toCreate.WalletID,
		Provider:       toCreate.Provider,
		ProviderUserID: toCreate.ProviderUserID,
		Credentials:    toCreate.Credentials,
	})
	return database.One(sqlcSeller, err, repo.dbe, mapSeller)
}

func (repo *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Seller, error) {
	sqlcSeller, err := repo.db(ctx).GetSellerByID(ctx, id)
	return database.One(sqlcSeller, err, repo.dbe, mapSeller)
}

func (repo *Repo) List(ctx context.Context) ([]models.Seller, error) {
	sqlcSellers, err := repo.db(ctx).ListSellers(ctx)
	return database.Many(sqlcSellers, err, repo.dbe, mapSeller)
}

func (repo *Repo) ListByWallet(ctx context.Context, walletID uuid.UUID) ([]models.Seller, error) {
	sqlcSellers, err := repo.db(ctx).ListSellersByWallet(ctx, walletID)
	return database.Many(sqlcSellers, err, repo.dbe, mapSeller)
}

func (repo *Repo) Revoke(ctx context.Context, id uuid.UUID) error {
	return repo.dbe(repo.db(ctx).RevokeSeller(ctx, id))
}
