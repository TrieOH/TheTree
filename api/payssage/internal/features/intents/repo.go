package intents

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

var _ ports.IntentRepo = (*Repo)(nil)

func NewRepo(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler(errs.ErrIntentNotFound),
	}
}

// db is the queries handle, on the context's transaction when it carries one.
func (repo *Repo) db(ctx context.Context) *sqlc.Queries { return database.Queries(ctx, repo.q) }

func mapIntent(src sqlc.Intent) models.Intent {
	var statusDetail *models.IntentStatusDetail
	if src.StatusDetail != nil {
		statusDetail = new(models.IntentStatusDetail(*src.StatusDetail))
	}

	return models.Intent{
		ID:            src.ID,
		WalletID:      src.WalletID,
		SellerID:      src.SellerID,
		CollectorID:   src.CollectorID,
		AmountCents:   src.AmountCents,
		Currency:      src.Currency,
		Sandbox:       src.Sandbox,
		Provider:      src.Provider,
		Status:        models.IntentStatus(src.Status),
		StatusDetail:  statusDetail,
		ProviderData:  src.ProviderData,
		Metadata:      src.Metadata,
		ExternalID:    src.ExternalID,
		ExternalGroup: src.ExternalGroup,
		CreatedAt:     src.CreatedAt,
		UpdatedAt:     src.UpdatedAt,
	}
}

func (repo *Repo) Create(ctx context.Context, toCreate models.Intent) (*models.Intent, error) {
	var statusDetail *string
	if toCreate.StatusDetail != nil {
		statusDetail = new(string(*toCreate.StatusDetail))
	}

	sqlcIntent, err := repo.db(ctx).CreateIntent(ctx, sqlc.CreateIntentParams{
		ID:            toCreate.ID,
		WalletID:      toCreate.WalletID,
		SellerID:      toCreate.SellerID,
		CollectorID:   toCreate.CollectorID,
		AmountCents:   toCreate.AmountCents,
		Currency:      toCreate.Currency,
		Sandbox:       toCreate.Sandbox,
		Provider:      toCreate.Provider,
		Status:        string(toCreate.Status),
		StatusDetail:  statusDetail,
		ProviderData:  toCreate.ProviderData,
		Metadata:      toCreate.Metadata,
		ExternalID:    toCreate.ExternalID,
		ExternalGroup: toCreate.ExternalGroup,
	})
	return database.One(sqlcIntent, err, repo.dbe, mapIntent)
}

func (repo *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Intent, error) {
	sqlcIntent, err := repo.db(ctx).GetIntentByID(ctx, id)
	return database.One(sqlcIntent, err, repo.dbe, mapIntent)
}

func (repo *Repo) GetByProviderTransactionID(ctx context.Context, provider string, transactionID string) (*models.Intent, error) {
	sqlcIntent, err := repo.db(ctx).GetIntentByProviderTransactionID(ctx, sqlc.GetIntentByProviderTransactionIDParams{
		Provider:      provider,
		TransactionID: transactionID,
	})
	return database.One(sqlcIntent, err, repo.dbe, mapIntent)
}

func (repo *Repo) ListByOrg(ctx context.Context, orgID uuid.UUID) ([]models.Intent, error) {
	sqlcIntents, err := repo.db(ctx).ListIntentsByOrg(ctx, &orgID)
	return database.Many(sqlcIntents, err, repo.dbe, mapIntent)
}

func (repo *Repo) ListByOwner(ctx context.Context, ownerID uuid.UUID) ([]models.Intent, error) {
	sqlcIntents, err := repo.db(ctx).ListIntentsByOwner(ctx, ownerID)
	return database.Many(sqlcIntents, err, repo.dbe, mapIntent)
}

func (repo *Repo) ListByWallet(ctx context.Context, walletID uuid.UUID) ([]models.Intent, error) {
	sqlcIntents, err := repo.db(ctx).ListIntentsByWallet(ctx, walletID)
	return database.Many(sqlcIntents, err, repo.dbe, mapIntent)
}

func (repo *Repo) Update(ctx context.Context, toUpdate models.Intent) (*models.Intent, error) {
	var statusDetail *string
	if toUpdate.StatusDetail != nil {
		statusDetail = new(string(*toUpdate.StatusDetail))
	}

	sqlcIntent, err := repo.db(ctx).UpdateIntent(ctx, sqlc.UpdateIntentParams{
		ID:           toUpdate.ID,
		Status:       string(toUpdate.Status),
		StatusDetail: statusDetail,
		ProviderData: toUpdate.ProviderData,
	})
	return database.One(sqlcIntent, err, repo.dbe, mapIntent)
}
