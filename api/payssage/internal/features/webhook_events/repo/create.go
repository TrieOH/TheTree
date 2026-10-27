package repo

import (
	"context"
	"lib/database"
	"payssage/internal/sqlc"
	"payssage/models"
)

func (r *Repo) Create(ctx context.Context, toCreate models.WebhookEvent) (*models.WebhookEvent, error) {
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
