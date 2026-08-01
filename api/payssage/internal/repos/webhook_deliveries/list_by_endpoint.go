package webhook_deliveries

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"lib/xslices"
	"payssage/models"

	"github.com/google/uuid"
)

func (r *Repo) ListByEndpoint(ctx context.Context, endpointID uuid.UUID) ([]models.WebhookDelivery, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListByEndpoint")
	defer span.End()
	rows, err := database.Queries(ctx, r.q).ListWebhookDeliveriesByEndpoint(ctx, endpointID)
	if err != nil {
		return nil, r.dbe(err)
	}
	return xslices.MapSlice(rows, mapWebhookDelivery), nil
}
