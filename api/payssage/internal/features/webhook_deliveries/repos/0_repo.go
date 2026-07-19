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

var _ ports.WebhookDeliveryRepo = (*repo)(nil)

func NewRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.WebhookDeliveryRepo {
	return &repo{
		q:      q,
		log:    log,
		tracer: tracer,
		dbe:    database.NewErrorHandler("webhook_delivery"),
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
