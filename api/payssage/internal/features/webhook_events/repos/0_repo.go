package repos

import (
	"lib/database"
	sqlc2 "payssage/internal/sqlc"
	"payssage/models"
	"payssage/ports"

	"go.opentelemetry.io/otel/trace"
)

type repo struct {
	q      *sqlc2.Queries
	tracer trace.Tracer
	dbe    database.ErrorHandler
}

var _ ports.WebhookEventRepo = (*repo)(nil)

func NewRepo(q *sqlc2.Queries, tracer trace.Tracer) ports.WebhookEventRepo {
	return &repo{
		q:      q,
		tracer: tracer,
		dbe:    database.NewErrorHandler("webhook_event"),
	}
}

func mapWebhookEvent(src sqlc2.WebhookEvent) models.WebhookEvent {
	var externalID string
	if src.ExternalID != nil {
		externalID = *src.ExternalID
	}
	return models.WebhookEvent{
		ID:         src.ID,
		WalletID:   src.WalletID,
		IntentID:   src.IntentID,
		Provider:   src.Provider,
		ExternalID: externalID,
		EventType:  src.EventType,
		Payload:    src.Payload,
		ReceivedAt: src.ReceivedAt,
	}
}
