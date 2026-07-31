package repos

import (
	"lib/database"
	sqlc2 "payssage/internal/sqlc"
	"payssage/models"
	"payssage/ports"
)

type Repo struct {
	q   *sqlc2.Queries
	dbe database.ErrorHandler
}

var _ ports.WebhookEndpointRepo = (*Repo)(nil)

func NewRepo(q *sqlc2.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler("webhook_endpoint"),
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
