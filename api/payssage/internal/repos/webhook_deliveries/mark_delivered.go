package webhook_deliveries

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"payssage/models"

	"github.com/google/uuid"
)

func (r *Repo) MarkDelivered(ctx context.Context, id uuid.UUID) (*models.WebhookDelivery, error) {
	ctx, span := telemetry.StartSpan(ctx, "MarkDelivered")
	defer span.End()
	row, err := database.Queries(ctx, r.q).MarkDeliveryDelivered(ctx, id)
	if err != nil {
		return nil, r.dbe(err)
	}
	return new(mapWebhookDelivery(row)), nil
}
