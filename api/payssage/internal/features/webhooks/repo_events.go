package webhooks

import (
	"context"

	"payssage/internal/platform/database"
	"payssage/internal/platform/database/sqlc"
	"payssage/internal/shared/contracts"
	"payssage/internal/shared/errx"
	"payssage/internal/shared/ports"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type webhookEventRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	tracer trace.Tracer
}

var _ ports.WebhookEventRepo = (*webhookEventRepo)(nil)

func NewWebhookEventRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.WebhookEventRepo {
	return &webhookEventRepo{q: q, log: log, tracer: tracer}
}

func (repo *webhookEventRepo) queries(ctx context.Context) *sqlc.Queries {
	if tx, ok := ctx.Value(database.TxKeyValue).(pgx.Tx); ok && tx != nil {
		return repo.q.WithTx(tx)
	}
	return repo.q
}

func mapWebhookEventFromDB(src *sqlc.WebhookEvent) *contracts.WebhookEventOriginal {
	return &contracts.WebhookEventOriginal{
		ID:          src.ID,
		WorkspaceID: src.WorkspaceID,
		IntentID:    src.IntentID,
		Provider:    src.Provider,
		ExternalID:  src.ExternalID,
		EventType:   src.EventType,
		Payload:     src.Payload,
		ReceivedAt:  src.ReceivedAt,
	}
}

func (repo *webhookEventRepo) Create(ctx context.Context, toCreate contracts.WebhookEventOriginal) (*contracts.WebhookEventOriginal, error) {
	ctx, span := repo.tracer.Start(ctx, "WebhookEventRepo.Create")
	defer span.End()

	row, err := repo.queries(ctx).CreateWebhookEvent(ctx, sqlc.CreateWebhookEventParams{
		ID:        toCreate.ID,
		Provider:  toCreate.Provider,
		EventType: toCreate.EventType,
		Payload:   toCreate.Payload,
	})
	if err != nil {
		return nil, errx.FromDB(err, "webhook_event")
	}

	return mapWebhookEventFromDB(&row), nil
}

func (repo *webhookEventRepo) Enrich(ctx context.Context, id, workspaceID, intentID uuid.UUID, externalID string) (*contracts.WebhookEventOriginal, error) {
	ctx, span := repo.tracer.Start(ctx, "WebhookEventRepo.Enrich")
	defer span.End()

	row, err := repo.queries(ctx).EnrichWebhookEvent(ctx, sqlc.EnrichWebhookEventParams{
		ID:          id,
		WorkspaceID: &workspaceID,
		IntentID:    &intentID,
		ExternalID:  &externalID,
	})
	if err != nil {
		return nil, errx.FromDB(err, "webhook_event")
	}

	return mapWebhookEventFromDB(&row), nil
}

func (repo *webhookEventRepo) GetByID(ctx context.Context, id uuid.UUID) (*contracts.WebhookEventOriginal, error) {
	ctx, span := repo.tracer.Start(ctx, "WebhookEventRepo.GetByID")
	defer span.End()

	row, err := repo.queries(ctx).GetWebhookEventByID(ctx, id)
	if err != nil {
		return nil, errx.FromDB(err, "webhook_event")
	}

	return mapWebhookEventFromDB(&row), nil
}

func (repo *webhookEventRepo) ListByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]contracts.WebhookEventOriginal, error) {
	ctx, span := repo.tracer.Start(ctx, "WebhookEventRepo.ListByWorkspace")
	defer span.End()

	rows, err := repo.queries(ctx).ListWebhookEventsByWorkspace(ctx, &workspaceID)
	if err != nil {
		return nil, errx.FromDB(err, "webhook_event")
	}

	out := make([]contracts.WebhookEventOriginal, 0, len(rows))
	for _, row := range rows {
		out = append(out, *mapWebhookEventFromDB(&row))
	}
	return out, nil
}

func (repo *webhookEventRepo) ListByProvider(ctx context.Context, provider string) ([]contracts.WebhookEventOriginal, error) {
	ctx, span := repo.tracer.Start(ctx, "WebhookEventRepo.ListByProvider")
	defer span.End()

	rows, err := repo.queries(ctx).ListWebhookEventsByProvider(ctx, provider)
	if err != nil {
		return nil, errx.FromDB(err, "webhook_event")
	}

	out := make([]contracts.WebhookEventOriginal, 0, len(rows))
	for _, row := range rows {
		out = append(out, *mapWebhookEventFromDB(&row))
	}
	return out, nil
}
