package repos

import (
	"lib/database"
	"payssage/internal/database/sqlc"
	"payssage/models"
	"payssage/ports"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type repo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	tracer trace.Tracer
	dbe    database.ErrorHandler
}

var _ ports.WebhookEndpointRepo = (*repo)(nil)

func NewRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.WebhookEndpointRepo {
	return &repo{
		q:      q,
		log:    log,
		tracer: tracer,
		dbe:    database.NewErrorHandler("webhook_endpoint"),
	}
}

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
