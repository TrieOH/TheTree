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

var _ ports.WebhookDeliveryRepo = (*Repo)(nil)

func New(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewNamedErrorHandler(errs.ErrWebhookDeliveryNotFound),
	}
}

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
