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

type webhookEndpointRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	tracer trace.Tracer
}

var _ ports.WebhookEndpointRepo = (*webhookEndpointRepo)(nil)

func NewWebhookEndpointRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.WebhookEndpointRepo {
	return &webhookEndpointRepo{q: q, log: log, tracer: tracer}
}

func (repo *webhookEndpointRepo) queries(ctx context.Context) *sqlc.Queries {
	if tx, ok := ctx.Value(database.TxKeyValue).(pgx.Tx); ok && tx != nil {
		return repo.q.WithTx(tx)
	}
	return repo.q
}

func mapWebhookEndpointFromDB(src *sqlc.WebhookEndpoint) *contracts.WebhookEndpoint {
	return &contracts.WebhookEndpoint{
		ID:          src.ID,
		ScopeID:     src.ScopeID,
		WorkspaceID: src.WorkspaceID,
		URL:         src.Url,
		Secret:      src.Secret,
		CreatedAt:   src.CreatedAt,
		DeletedAt:   src.DeletedAt,
	}
}

func (repo *webhookEndpointRepo) Create(ctx context.Context, toCreate contracts.WebhookEndpoint) (*contracts.WebhookEndpoint, error) {
	ctx, span := repo.tracer.Start(ctx, "WebhookEndpointRepo.Create")
	defer span.End()

	row, err := repo.queries(ctx).CreateWebhookEndpoint(ctx, sqlc.CreateWebhookEndpointParams{
		ID:          toCreate.ID,
		ScopeID:     toCreate.ScopeID,
		WorkspaceID: toCreate.WorkspaceID,
		Url:         toCreate.URL,
		Secret:      toCreate.Secret,
	})
	if err != nil {
		return nil, errx.FromDB(err, "webhook_endpoint")
	}

	return mapWebhookEndpointFromDB(&row), nil
}

func (repo *webhookEndpointRepo) GetByID(ctx context.Context, id uuid.UUID) (*contracts.WebhookEndpoint, error) {
	ctx, span := repo.tracer.Start(ctx, "WebhookEndpointRepo.GetByID")
	defer span.End()

	row, err := repo.queries(ctx).GetWebhookEndpointByID(ctx, id)
	if err != nil {
		return nil, errx.FromDB(err, "webhook_endpoint")
	}

	return mapWebhookEndpointFromDB(&row), nil
}

func (repo *webhookEndpointRepo) ListByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]contracts.WebhookEndpoint, error) {
	ctx, span := repo.tracer.Start(ctx, "WebhookEndpointRepo.ListByWorkspace")
	defer span.End()

	rows, err := repo.queries(ctx).ListWebhookEndpointsByWorkspace(ctx, workspaceID)
	if err != nil {
		return nil, errx.FromDB(err, "webhook_endpoint")
	}

	out := make([]contracts.WebhookEndpoint, 0, len(rows))
	for _, row := range rows {
		out = append(out, *mapWebhookEndpointFromDB(&row))
	}
	return out, nil
}

func (repo *webhookEndpointRepo) Delete(ctx context.Context, id, workspaceID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "WebhookEndpointRepo.Delete")
	defer span.End()

	err := repo.queries(ctx).DeleteWebhookEndpoint(ctx, sqlc.DeleteWebhookEndpointParams{
		ID:          id,
		WorkspaceID: workspaceID,
	})
	if err != nil {
		return errx.FromDB(err, "webhook_endpoint")
	}

	return nil
}
