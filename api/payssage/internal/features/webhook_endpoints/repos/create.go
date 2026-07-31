package repos

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"payssage/internal/sqlc"
	"payssage/models"
)

func (r *Repo) Create(ctx context.Context, toCreate models.WebhookEndpoint) (*models.WebhookEndpoint, error) {
	ctx, span := telemetry.StartSpan(ctx, "Create")
	defer span.End()
	row, err := database.Queries(ctx, r.q).CreateWebhookEndpoint(ctx, sqlc.CreateWebhookEndpointParams{
		WalletID: toCreate.WalletID,
		Name:     toCreate.Name,
		Url:      toCreate.URL,
		Secret:   toCreate.Secret,
	})
	if err != nil {
		return nil, r.dbe(err)
	}
	return new(mapWebhookEndpoint(row)), nil
}
