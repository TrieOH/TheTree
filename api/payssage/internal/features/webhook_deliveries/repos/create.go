package repos

import (
	"context"
	"lib/database"
	"payssage/internal/database/sqlc"
	"payssage/models"
)

func (r *repo) Create(ctx context.Context, toCreate models.WebhookDelivery) (*models.WebhookDelivery, error) {
	ctx, span := r.tracer.Start(ctx, "Create")
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
