package repo

import (
	"lib/database"
	"payssage/errs"
	"payssage/internal/sqlc"
	"payssage/models"
	"payssage/ports"
)

type Repo struct {
	q   *sqlc.Queries
	dbe database.NamedErrorHandler
}

var _ ports.WebhookEndpointRepo = (*Repo)(nil)

func New(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewNamedErrorHandler(errs.ErrWebhookEndpointNotFound),
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
