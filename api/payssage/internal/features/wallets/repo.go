package wallets

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

var _ ports.WalletRepo = (*Repo)(nil)

func NewRepo(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler(errs.ErrWalletNotFound),
	}
}

// db is the queries handle, on the context's transaction when it carries one.
func (repo *Repo) db(ctx context.Context) *sqlc.Queries { return database.Queries(ctx, repo.q) }

func mapWallet(src sqlc.Wallet) models.Wallet {
	return models.Wallet{
		ID:             src.ID,
		OwnerID:        src.OwnerID,
		OrganizationID: src.OrganizationID,
		Name:           src.Name,
		Sandbox:        src.Sandbox,
		FeeBps:         src.FeeBps,
		CollectorID:    src.CollectorID,
		CreatedAt:      src.CreatedAt,
	}
}

func (repo *Repo) BindCollector(ctx context.Context, walletID, collectorID uuid.UUID) error {
	return repo.dbe(repo.db(ctx).BindCollector(ctx, sqlc.BindCollectorParams{
		ID:          walletID,
		CollectorID: &collectorID,
	}))
}

func (repo *Repo) Create(ctx context.Context, toCreate models.Wallet) (*models.Wallet, error) {
	sqlcWallet, err := repo.db(ctx).CreateWallet(ctx, sqlc.CreateWalletParams{
		OwnerID:        toCreate.OwnerID,
		OrganizationID: toCreate.OrganizationID,
		Name:           toCreate.Name,
		Sandbox:        toCreate.Sandbox,
		FeeBps:         toCreate.FeeBps,
	})
	return database.One(sqlcWallet, err, repo.dbe, mapWallet)
}

func (repo *Repo) GetByID(ctx context.Context, walletID uuid.UUID) (*models.Wallet, error) {
	sqlcWallet, err := repo.db(ctx).GetWalletByID(ctx, walletID)
	return database.One(sqlcWallet, err, repo.dbe, mapWallet)
}

func (repo *Repo) GetRole(ctx context.Context, actorID, walletID uuid.UUID) (models.OrganizationRole, error) {
	role, err := repo.db(ctx).GetWalletRole(ctx, sqlc.GetWalletRoleParams{
		ID:      walletID,
		ActorID: actorID,
	})
	if err != nil {
		return "", repo.dbe(err)
	}
	return models.OrganizationRole(role), nil
}

func (repo *Repo) List(ctx context.Context, ownerID uuid.UUID) ([]models.Wallet, error) {
	sqlcWallets, err := repo.db(ctx).ListWallets(ctx, ownerID)
	return database.Many(sqlcWallets, err, repo.dbe, mapWallet)
}

func (repo *Repo) ListFromOrg(ctx context.Context, orgID uuid.UUID) ([]models.Wallet, error) {
	sqlcWallets, err := repo.db(ctx).ListOrgWallets(ctx, &orgID)
	return database.Many(sqlcWallets, err, repo.dbe, mapWallet)
}

func (repo *Repo) SetFeeBPS(ctx context.Context, walletID uuid.UUID, feeBPS int) error {
	err := repo.db(ctx).SetWalletFeeBPS(ctx, sqlc.SetWalletFeeBPSParams{
		FeeBps: feeBPS,
		ID:     walletID,
	})
	return repo.dbe(err)
}

func (repo *Repo) SetSandboxState(ctx context.Context, walletID uuid.UUID, state bool) error {
	err := repo.db(ctx).SetWalletSandboxState(ctx, sqlc.SetWalletSandboxStateParams{
		Sandbox: state,
		ID:      walletID,
	})
	return repo.dbe(err)
}

func (repo *Repo) UnbindCollector(ctx context.Context, walletID uuid.UUID) error {
	return repo.dbe(repo.db(ctx).UnbindCollector(ctx, walletID))
}
