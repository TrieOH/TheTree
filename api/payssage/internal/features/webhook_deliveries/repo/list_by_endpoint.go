package repo

import (
	"context"
	"lib/database"
	"lib/xslices"
	"payssage/models"

	"github.com/google/uuid"
)

func (r *Repo) ListByEndpoint(ctx context.Context, endpointID uuid.UUID) ([]models.WebhookDelivery, error) {
	rows, err := database.Queries(ctx, r.q).ListWebhookDeliveriesByEndpoint(ctx, endpointID)
	if err != nil {
		return nil, r.dbe(err)
	}
	return xslices.MapSlice(rows, mapWebhookDelivery), nil
}
