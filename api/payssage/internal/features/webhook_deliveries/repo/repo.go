package repo

import (
	"context"
	"lib/database"
	"payssage/errs"
	"payssage/internal/sqlc"
	"payssage/models"
	"payssage/ports"
	"time"

	"github.com/google/uuid"
)

type Repo struct {
	q   *sqlc.Queries
	dbe database.ErrorHandler
}

var _ ports.WebhookDeliveryRepo = (*Repo)(nil)

func New(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler(errs.ErrWebhookDeliveryNotFound),
	}
}

// db is the queries handle, on the context's transaction when it carries one.
func (r *Repo) db(ctx context.Context) *sqlc.Queries { return database.Queries(ctx, r.q) }

func mapWebhookDelivery(src sqlc.WebhookDelivery) models.WebhookDelivery {
	return models.WebhookDelivery{
		ID:              src.ID,
		EndpointID:      src.EndpointID,
		EventID:         src.EventID,
		Status:          models.WebhookDeliveryStatus(src.Status),
		Attempts:        src.Attempts,
		LastAttemptedAt: src.LastAttemptedAt,
		ResponseStatus:  src.ResponseStatus,
		ResponseBody:    src.ResponseBody,
		CreatedAt:       src.CreatedAt,
		UpdatedAt:       src.UpdatedAt,
	}
}

func (r *Repo) Create(ctx context.Context, toCreate models.WebhookDelivery) (*models.WebhookDelivery, error) {
	row, err := r.db(ctx).CreateWebhookDelivery(ctx, sqlc.CreateWebhookDeliveryParams{
		EndpointID: toCreate.EndpointID,
		EventID:    toCreate.EventID,
		Status:     string(toCreate.Status),
	})
	return database.One(row, err, r.dbe, mapWebhookDelivery)
}

func (r *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.WebhookDelivery, error) {
	row, err := r.db(ctx).GetWebhookDeliveryByID(ctx, id)
	return database.One(row, err, r.dbe, mapWebhookDelivery)
}

func (r *Repo) IncrementAttempt(ctx context.Context, id uuid.UUID) (*models.WebhookDelivery, error) {
	row, err := r.db(ctx).IncrementDeliveryAttempt(ctx, id)
	return database.One(row, err, r.dbe, mapWebhookDelivery)
}

func (r *Repo) ListByEndpoint(ctx context.Context, endpointID uuid.UUID) ([]models.WebhookDelivery, error) {
	rows, err := r.db(ctx).ListWebhookDeliveriesByEndpoint(ctx, endpointID)
	return database.Many(rows, err, r.dbe, mapWebhookDelivery)
}

func (r *Repo) ListUnsent(ctx context.Context, olderThan time.Duration, limit int) ([]models.WebhookDelivery, error) {
	rows, err := r.db(ctx).ListUnsentWebhookDeliveries(ctx, sqlc.ListUnsentWebhookDeliveriesParams{
		CreatedBefore: time.Now().Add(-olderThan),
		MaxRows:       limit,
	})
	return database.Many(rows, err, r.dbe, mapWebhookDelivery)
}

func (r *Repo) MarkDelivered(ctx context.Context, id uuid.UUID) (*models.WebhookDelivery, error) {
	row, err := r.db(ctx).MarkDeliveryDelivered(ctx, id)
	return database.One(row, err, r.dbe, mapWebhookDelivery)
}

func (r *Repo) MarkFailed(ctx context.Context, id uuid.UUID) (*models.WebhookDelivery, error) {
	row, err := r.db(ctx).MarkDeliveryFailed(ctx, id)
	return database.One(row, err, r.dbe, mapWebhookDelivery)
}

func (r *Repo) Update(ctx context.Context, params models.UpdateDeliveryParams) (*models.WebhookDelivery, error) {
	row, err := r.db(ctx).UpdateWebhookDelivery(ctx, sqlc.UpdateWebhookDeliveryParams{
		Status:          string(params.Status),
		Attempts:        params.Attempts,
		LastAttemptedAt: params.LastAttemptedAt,
		ResponseStatus:  params.ResponseStatus,
		ResponseBody:    params.ResponseBody,
		ID:              params.ID,
	})
	return database.One(row, err, r.dbe, mapWebhookDelivery)
}
