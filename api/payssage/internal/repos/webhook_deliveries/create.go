package webhook_deliveries

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"payssage/internal/sqlc"
	"payssage/models"
)

func (r *Repo) Create(ctx context.Context, toCreate models.WebhookDelivery) (*models.WebhookDelivery, error) {
	ctx, span := telemetry.StartSpan(ctx, "Create")
	defer span.End()
	row, err := database.Queries(ctx, r.q).CreateWebhookDelivery(ctx, sqlc.CreateWebhookDeliveryParams{
		EndpointID: toCreate.EndpointID,
		EventID:    toCreate.EventID,
		Status:     string(toCreate.Status),
	})
	if err != nil {
		return nil, r.dbe(err)
	}
	return new(mapWebhookDelivery(row)), nil
}
