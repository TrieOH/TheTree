package webhook_endpoints

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

var _ ports.WebhookEndpointRepo = (*Repo)(nil)

func NewRepo(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler(errs.ErrWebhookEndpointNotFound),
	}
}

// db is the queries handle, on the context's transaction when it carries one.
func (r *Repo) db(ctx context.Context) *sqlc.Queries { return database.Queries(ctx, r.q) }

func mapWebhookEndpoint(src sqlc.WebhookEndpoint) models.WebhookEndpoint {
	return models.WebhookEndpoint{
		ID:        src.ID,
		WalletID:  src.WalletID,
		Name:      src.Name,
		URL:       src.Url,
		Secret:    src.Secret,
		CreatedAt: src.CreatedAt,
	}
}

func (r *Repo) Create(ctx context.Context, toCreate models.WebhookEndpoint) (*models.WebhookEndpoint, error) {
	row, err := r.db(ctx).CreateWebhookEndpoint(ctx, sqlc.CreateWebhookEndpointParams{
		WalletID: toCreate.WalletID,
		Name:     toCreate.Name,
		Url:      toCreate.URL,
		Secret:   toCreate.Secret,
	})
	return database.One(row, err, r.dbe, mapWebhookEndpoint)
}

func (r *Repo) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.db(ctx).DeleteWebhookEndpoint(ctx, id)
	return r.dbe(err)
}

func (r *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.WebhookEndpoint, error) {
	row, err := r.db(ctx).GetWebhookEndpointByID(ctx, id)
	return database.One(row, err, r.dbe, mapWebhookEndpoint)
}

func (r *Repo) ListByWallet(ctx context.Context, walletID uuid.UUID) ([]models.WebhookEndpoint, error) {
	rows, err := r.db(ctx).ListWebhookEndpointsByWallet(ctx, walletID)
	return database.Many(rows, err, r.dbe, mapWebhookEndpoint)
}
