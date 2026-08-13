package webhook_events

import (
	"lib/database"
	"payssage/internal/sqlc"
	"payssage/models"
	"payssage/ports"
)

type Repo struct {
	q   *sqlc.Queries
	dbe database.ErrorHandler
}

var _ ports.WebhookEventRepo = (*Repo)(nil)

func NewRepo(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler("webhook_event"),
	}
}

func mapWebhookEvent(src sqlc.WebhookEvent) models.WebhookEvent {
	var externalID string
	if src.ExternalID != nil {
		externalID = *src.ExternalID
	}
	return models.WebhookEvent{
		ID:           src.ID,
		WalletID:     src.WalletID,
		IntentID:     src.IntentID,
		Provider:     src.Provider,
		ExternalID:   externalID,
		EventType:    src.EventType,
		StatusDetail: src.StatusDetail,
		Payload:      src.Payload,
		ReceivedAt:   src.ReceivedAt,
	}
}
