package infrastructure

import (
	"TriePayments/internal/core/domain"
	"TriePayments/internal/plataform/database"
	"TriePayments/internal/plataform/database/sqlc"
	"TriePayments/internal/shared/errx"
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type webhookDeliveryRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	tracer trace.Tracer
}

var _ domain.WebhookDeliveryRepo = (*webhookDeliveryRepo)(nil)

func NewWebhookDeliveryRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) domain.WebhookDeliveryRepo {
	return &webhookDeliveryRepo{q: q, log: log, tracer: tracer}
}

func (repo *webhookDeliveryRepo) queries(ctx context.Context) *sqlc.Queries {
	if tx, ok := ctx.Value(database.TxKeyValue).(pgx.Tx); ok && tx != nil {
		return repo.q.WithTx(tx)
	}
	return repo.q
}

func mapWebhookDeliveryFromDB(src *sqlc.WebhookDelivery) *domain.WebhookDelivery {
	return &domain.WebhookDelivery{
		ID:              src.ID,
		EndpointID:      src.EndpointID,
		IntentID:        src.IntentID,
		Event:           src.Event,
		Payload:         json.RawMessage(src.Payload),
		Status:          domain.DeliveryStatus(src.Status),
		Attempts:        int(src.Attempts),
		LastAttemptedAt: src.LastAttemptedAt,
		CreatedAt:       src.CreatedAt,
	}
}

func (repo *webhookDeliveryRepo) Create(ctx context.Context, toCreate domain.WebhookDelivery) (*domain.WebhookDelivery, error) {
	ctx, span := repo.tracer.Start(ctx, "WebhookDeliveryRepo.Create")
	defer span.End()

	row, err := repo.queries(ctx).CreateWebhookDelivery(ctx, sqlc.CreateWebhookDeliveryParams{
		ID:         toCreate.ID,
		EndpointID: toCreate.EndpointID,
		IntentID:   toCreate.IntentID,
		Event:      toCreate.Event,
		Payload:    toCreate.Payload,
		Status:     sqlc.DeliveryStatus(toCreate.Status),
	})
	if err != nil {
		return nil, errx.FromDB(err, "webhook_delivery")
	}

	return mapWebhookDeliveryFromDB(&row), nil
}

func (repo *webhookDeliveryRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.WebhookDelivery, error) {
	ctx, span := repo.tracer.Start(ctx, "WebhookDeliveryRepo.GetByID")
	defer span.End()

	row, err := repo.queries(ctx).GetWebhookDeliveryByID(ctx, id)
	if err != nil {
		return nil, errx.FromDB(err, "webhook_delivery")
	}

	return mapWebhookDeliveryFromDB(&row), nil
}

func (repo *webhookDeliveryRepo) ListByEndpoint(ctx context.Context, endpointID uuid.UUID) ([]domain.WebhookDelivery, error) {
	ctx, span := repo.tracer.Start(ctx, "WebhookDeliveryRepo.ListByEndpoint")
	defer span.End()

	rows, err := repo.queries(ctx).ListWebhookDeliveriesByEndpoint(ctx, endpointID)
	if err != nil {
		return nil, errx.FromDB(err, "webhook_delivery")
	}

	out := make([]domain.WebhookDelivery, 0, len(rows))
	for _, row := range rows {
		out = append(out, *mapWebhookDeliveryFromDB(&row))
	}
	return out, nil
}

func (repo *webhookDeliveryRepo) MarkDelivered(ctx context.Context, id uuid.UUID) (*domain.WebhookDelivery, error) {
	ctx, span := repo.tracer.Start(ctx, "WebhookDeliveryRepo.MarkDelivered")
	defer span.End()

	row, err := repo.queries(ctx).MarkDeliveryDelivered(ctx, id)
	if err != nil {
		return nil, errx.FromDB(err, "webhook_delivery")
	}

	return mapWebhookDeliveryFromDB(&row), nil
}

func (repo *webhookDeliveryRepo) MarkFailed(ctx context.Context, id uuid.UUID) (*domain.WebhookDelivery, error) {
	ctx, span := repo.tracer.Start(ctx, "WebhookDeliveryRepo.MarkFailed")
	defer span.End()

	row, err := repo.queries(ctx).MarkDeliveryFailed(ctx, id)
	if err != nil {
		return nil, errx.FromDB(err, "webhook_delivery")
	}

	return mapWebhookDeliveryFromDB(&row), nil
}

func (repo *webhookDeliveryRepo) IncrementAttempt(ctx context.Context, id uuid.UUID) (*domain.WebhookDelivery, error) {
	ctx, span := repo.tracer.Start(ctx, "WebhookDeliveryRepo.IncrementAttempt")
	defer span.End()

	row, err := repo.queries(ctx).IncrementDeliveryAttempt(ctx, id)
	if err != nil {
		return nil, errx.FromDB(err, "webhook_delivery")
	}

	return mapWebhookDeliveryFromDB(&row), nil
}
