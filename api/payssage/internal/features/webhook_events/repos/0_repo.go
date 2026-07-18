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

var _ ports.WebhookEventRepo = (*repo)(nil)

func NewRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.WebhookEventRepo {
	return &repo{
		q:      q,
		log:    log,
		tracer: tracer,
		dbe:    database.NewErrorHandler("webhook_event"),
	}
}

func mapWebhookEvent(src sqlc.WebhookEvent) models.WebhookEvent {
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
