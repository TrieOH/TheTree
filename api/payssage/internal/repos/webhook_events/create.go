package webhook_events

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"payssage/internal/sqlc"
	"payssage/models"
)

func (r *Repo) Create(ctx context.Context, toCreate models.WebhookEvent) (*models.WebhookEvent, error) {
	ctx, span := telemetry.StartSpan(ctx, "Create")
	defer span.End()
	row, err := database.Queries(ctx, r.q).CreateWebhookEvent(ctx, sqlc.CreateWebhookEventParams{
		WalletID:     toCreate.WalletID,
		IntentID:     toCreate.IntentID,
		Provider:     toCreate.Provider,
		ExternalID:   &toCreate.ExternalID,
		EventType:    toCreate.EventType,
		StatusDetail: toCreate.StatusDetail,
		Payload:      toCreate.Payload,
	})
	if err != nil {
		return nil, r.dbe(err)
	}
	return new(mapWebhookEvent(row)), nil
}
