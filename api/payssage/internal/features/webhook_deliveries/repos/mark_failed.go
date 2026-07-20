package repos

import (
	"context"
	"lib/database"
	"payssage/models"

	"github.com/google/uuid"
)

func (r *repo) MarkFailed(ctx context.Context, id uuid.UUID) (*models.WebhookDelivery, error) {
	ctx, span := r.tracer.Start(ctx, "MarkFailed")
	defer span.End()
	row, err := database.Queries(ctx, r.q).MarkDeliveryFailed(ctx, id)
	if err != nil {
		return nil, r.dbe(err)
	}
	return new(mapWebhookDelivery(row)), nil
}
