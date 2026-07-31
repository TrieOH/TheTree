package repos

import (
	"lib/database"
	sqlc2 "payssage/internal/sqlc"
	"payssage/models"
	"payssage/ports"

	"go.opentelemetry.io/otel/trace"
)

type Repo struct {
	q      *sqlc2.Queries
	tracer trace.Tracer
	dbe    database.ErrorHandler
}

var _ ports.WebhookEndpointRepo = (*Repo)(nil)

func NewRepo(q *sqlc2.Queries, tracer trace.Tracer) *Repo {
	return &Repo{
		q:      q,
		tracer: tracer,
		dbe:    database.NewErrorHandler("webhook_endpoint"),
	}
}

func mapWebhookEndpoint(src sqlc2.WebhookEndpoint) models.WebhookEndpoint {
	return models.WebhookEndpoint{
		ID:        src.ID,
		WalletID:  src.WalletID,
		Name:      src.Name,
		URL:       src.Url,
		Secret:    src.Secret,
		CreatedAt: src.CreatedAt,
	}
}
