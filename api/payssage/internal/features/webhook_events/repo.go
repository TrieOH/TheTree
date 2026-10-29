package webhook_events

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

var _ ports.WebhookEventRepo = (*Repo)(nil)

func NewRepo(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler(errs.ErrWebhookEventNotFound),
	}
}

// db is the queries handle, on the context's transaction when it carries one.
func (r *Repo) db(ctx context.Context) *sqlc.Queries { return database.Queries(ctx, r.q) }

func mapWebhookEvent(src sqlc.WebhookEvent) models.WebhookEvent {
	var externalID string
	if src.ExternalID != nil {
		externalID = *src.ExternalID
	}
	return models.WebhookEvent{
		ID:           src.ID,
		WalletID:     src.WalletID,
		IntentID:     src.IntentID,
		Provider:     src.Provider,
		ExternalID:   externalID,
		EventType:    src.EventType,
		StatusDetail: src.StatusDetail,
		Payload:      src.Payload,
		ReceivedAt:   src.ReceivedAt,
	}
}

func (r *Repo) Create(ctx context.Context, toCreate models.WebhookEvent) (*models.WebhookEvent, error) {
	row, err := r.db(ctx).CreateWebhookEvent(ctx, sqlc.CreateWebhookEventParams{
		WalletID:     toCreate.WalletID,
		IntentID:     toCreate.IntentID,
		Provider:     toCreate.Provider,
		ExternalID:   &toCreate.ExternalID,
		EventType:    toCreate.EventType,
		StatusDetail: toCreate.StatusDetail,
		Payload:      toCreate.Payload,
	})
	return database.One(row, err, r.dbe, mapWebhookEvent)
}

func (r *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.WebhookEvent, error) {
	row, err := r.db(ctx).GetWebhookEventByID(ctx, id)
	return database.One(row, err, r.dbe, mapWebhookEvent)
}

func (r *Repo) ListByProvider(ctx context.Context, provider string) ([]models.WebhookEvent, error) {
	rows, err := r.db(ctx).ListWebhookEventsByProvider(ctx, provider)
	return database.Many(rows, err, r.dbe, mapWebhookEvent)
}

func (r *Repo) ListByWallet(ctx context.Context, walletID uuid.UUID) ([]models.WebhookEvent, error) {
	rows, err := r.db(ctx).ListWebhookEventsByWallet(ctx, walletID)
	return database.Many(rows, err, r.dbe, mapWebhookEvent)
}
