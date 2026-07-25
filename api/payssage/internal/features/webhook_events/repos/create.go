package repos

import (
	"context"
	"lib/database"
	"payssage/internal/sqlc"
	"payssage/models"
)

func (r *repo) Create(ctx context.Context, toCreate models.WebhookEvent) (*models.WebhookEvent, error) {
	ctx, span := r.tracer.Start(ctx, "Create")
	defer span.End()
	row, err := database.Queries(ctx, r.q).CreateWebhookEvent(ctx, sqlc.CreateWebhookEventParams{
		WalletID:   toCreate.WalletID,
		IntentID:   toCreate.IntentID,
		Provider:   toCreate.Provider,
		ExternalID: &toCreate.ExternalID,
		EventType:  toCreate.EventType,
		Payload:    toCreate.Payload,
	})
	if err != nil {
		return nil, r.dbe(err)
	}
	return new(mapWebhookEvent(row)), nil
}
